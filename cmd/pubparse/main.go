package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ashahide/pubparse/internal/fileIO"
	"github.com/ashahide/pubparse/internal/jsonTools"
	"github.com/ashahide/pubparse/internal/makeReports"
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

	// Check 1-to-1 mapping
	if len(args.InputPath.Files) != len(args.OutputPath.Files) {
		return fmt.Errorf("input/output file count mismatch")
	}

	// Make a report.txt file in the output directory
	reportFile, err := filepath.Abs(filepath.Join(args.OutputPath.Path, "report.tsv"))
	if err != nil {
		return fmt.Errorf("failed to determine absolute path for report file: %w", err)
	}
	if err := fileIO.MakeFile(reportFile); err != nil {
		return fmt.Errorf("failed to create report file %q: %w", reportFile, err)
	}
	report, err := os.OpenFile(reportFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open report file %q: %w", reportFile, err)
	}

	// Write to report file
	startTime := time.Now()
	if _, err := report.WriteString(fmt.Sprintf("\n>>> Starting Time: %s", startTime)); err != nil {
		return fmt.Errorf("failed to write to report file %q: %w", reportFile, err)
	}

	// Write to report file
	if _, err := report.WriteString(fmt.Sprintf("\n>>> Input Directory: %s", args.InputPath.Path)); err != nil {
		return fmt.Errorf("failed to write to report file %q: %w", reportFile, err)
	}

	// Write to report file
	if _, err := report.WriteString(fmt.Sprintf("\n>>> Output Directory: %s", args.OutputPath.Path)); err != nil {
		return fmt.Errorf("failed to write to report file %q: %w", reportFile, err)
	}

	// Write to report file
	if _, err := report.WriteString(fmt.Sprintf("\n>>> Number of Inputs: %d", len(args.InputPath.Files))); err != nil {
		return fmt.Errorf("failed to write to report file %q: %w", reportFile, err)
	}

	// Print Inputs and Outputs
	fmt.Println(">>> Input Path:", args.InputPath.Path)
	fmt.Println(">>> Output Path:", args.OutputPath.Path)
	fmt.Println(">>> Number of Inputs:", len(args.InputPath.Files))
	fmt.Println(">>> Starting Time:", startTime)

	// Process each file
	for i := range args.InputPath.Files {

		makeReports.PrintProgressBar(i+1, len(args.InputPath.Files), startTime)

		fin := args.InputPath.Files[i]
		fout := args.OutputPath.Files[i]

		if err := fileIO.MakeFile(fout); err != nil {
			return fmt.Errorf("failed to create output file %q: %w", fout, err)
		}

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
			xmlTools.NormalizePubmedArticleSet(v)
			// Build the path to the validation schema.
			schemaPath := filepath.Join("internal", "jsonTools", "pubmed_json_schema.json")

			if err := jsonTools.ConvertToJson(v, fout, schemaPath); err != nil {
				return fmt.Errorf("failed to convert PubMed article to JSON %q: %w", fout, err)
			}
		case *xmlTools.PubmedBookArticleSet:
			xmlTools.NormalizePubmedArticleSet(v)
			schemaPath := filepath.Join("internal", "jsonTools", "pubmed_json_schema.json")
			if err := jsonTools.ConvertToJson(v, fout, schemaPath); err != nil {
				return fmt.Errorf("failed to convert PubMed book to JSON %q: %w", fout, err)
			}
		case *xmlTools.PMCArticle:
			xmlTools.NormalizePMCArticle(v)
			schemaPath := filepath.Join("internal", "jsonTools", "pmc_json_schema.json")
			if err := jsonTools.ConvertToJson(v, fout, schemaPath); err != nil {
				return fmt.Errorf("failed to convert PMC article to JSON %q: %w", fout, err)
			}
		default:
			return fmt.Errorf("unsupported data type for file %q", fin)
		}

		// Write to report file
		if _, err := report.WriteString(fmt.Sprintf("\n>>> Input file: %s\t Output file: %s\n", fin, fout)); err != nil {
			return fmt.Errorf("failed to write to report file %q: %w", reportFile, err)
		}
		if err := report.Sync(); err != nil {
			return fmt.Errorf("failed to sync report file %q: %w", reportFile, err)
		}

	}

	// Close the report file
	if err := report.Close(); err != nil {
		return fmt.Errorf("failed to close report file %q: %w", reportFile, err)
	}
	// Print completion message
	fmt.Println("\n>>> Finished processing files.")
	fmt.Println(">>> Report file:", reportFile)
	fmt.Println(">>> Elapsed Time:", time.Since(startTime))
	fmt.Println(">>> Exiting...")

	return nil
}
