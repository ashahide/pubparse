package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ashahide/pubparse/internal/fileIO"
	"github.com/ashahide/pubparse/internal/jsonTools"
	"github.com/ashahide/pubparse/internal/xmlTools"
)

func main() {
	if err := run(); err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pubparse [pubmed|pmc] -i input_path -o output_path")
		os.Exit(1)
	}

	mode := os.Args[1]
	var args fileIO.Arguments

	switch mode {
	case "pubmed":
		cmd := flag.NewFlagSet(mode, flag.ExitOnError)
		cmd.StringVar(&args.InputPath.Path, "i", "", "Path to the input file or directory")
		cmd.StringVar(&args.OutputPath.Path, "o", "", "Path to the output file or directory")
		if err := cmd.Parse(os.Args[2:]); err != nil {
			return err
		}

	case "pmc":
		cmd := flag.NewFlagSet(mode, flag.ExitOnError)
		cmd.StringVar(&args.InputPath.Path, "i", "", "Path to the input file or directory")
		cmd.StringVar(&args.OutputPath.Path, "o", "", "Path to the output file or directory")
		if err := cmd.Parse(os.Args[2:]); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown subcommand: %s\nUsage: pubparse [pubmed|pmc] -i input -o output", mode)
	}

	// Validate and gather input/output files
	if err := fileIO.HandleInputs(&args); err != nil {
		return fmt.Errorf("input handling failed: %w", err)
	}
	if err := fileIO.HandleOutputs(&args); err != nil {
		return fmt.Errorf("output handling failed: %w", err)
	}

	// Show file mappings
	fmt.Println("\nInput Path:", args.InputPath.Path)
	for _, f := range args.InputPath.Files {
		fmt.Printf("- %s\n", f)
	}
	fmt.Println("\nOutput Path:", args.OutputPath.Path)
	for _, f := range args.OutputPath.Files {
		fmt.Printf("- %s\n", f)
	}

	// Check 1-to-1 mapping
	if len(args.InputPath.Files) != len(args.OutputPath.Files) {
		return fmt.Errorf("input/output file count mismatch")
	}

	// Process each file
	for i := range args.InputPath.Files {
		fin := args.InputPath.Files[i]
		fout := args.OutputPath.Files[i]

		if err := fileIO.MakeFile(fout); err != nil {
			return fmt.Errorf("failed to create output file %q: %w", fout, err)
		}

		fmt.Println("\nProcessing file:", fin)

		var data interface{}
		var xml_parse_err error
		switch mode {
		case "pubmed", "pmc":
			data, xml_parse_err = xmlTools.ParsePubmedXML(fin)
		}

		if xml_parse_err != nil {
			return fmt.Errorf("failed to parse XML %q: %w", fin, xml_parse_err)
		}

		switch v := data.(type) {
		case *xmlTools.PubmedArticleSet:
			fmt.Println("Detected: PubmedArticleSet")
			xmlTools.NormalizePubmedArticleSet(v)
			// Build the path to the validation schema.
			schemaPath := filepath.Join("internal", "jsonTools", "pubmed_json_schema.json")

			if err := jsonTools.ConvertToJson(v, fout, schemaPath); err != nil {
				return fmt.Errorf("failed to convert PubMed article to JSON %q: %w", fout, err)
			}
		case *xmlTools.PubmedBookArticleSet:
			xmlTools.NormalizePubmedArticleSet(v)
			fmt.Println("Detected: PubmedBookArticleSet")
			schemaPath := filepath.Join("internal", "jsonTools", "pubmed_json_schema.json")
			if err := jsonTools.ConvertToJson(v, fout, schemaPath); err != nil {
				return fmt.Errorf("failed to convert PubMed book to JSON %q: %w", fout, err)
			}
		case *xmlTools.PMCArticle:
			fmt.Println("Detected: PMCArticle")
			xmlTools.NormalizePMCArticle(v)
			schemaPath := filepath.Join("internal", "jsonTools", "pmc_json_schema.json")
			if err := jsonTools.ConvertToJson(v, fout, schemaPath); err != nil {
				return fmt.Errorf("failed to convert PMC article to JSON %q: %w", fout, err)
			}
		default:
			return fmt.Errorf("unsupported data type for file %q", fin)
		}

	}

	return nil
}
