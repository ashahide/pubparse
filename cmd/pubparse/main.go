package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	var args fileIO.Arguments

	flag.StringVar(&args.InputPath.Path, "i", "", "Path to the input file or directory of files")
	flag.Parse()

	// Handle inputs and outputs
	if err := fileIO.HandleInputs(&args); err != nil {
		return fmt.Errorf("input handling failed: %w", err)
	}
	if err := fileIO.HandleOutputs(&args); err != nil {
		return fmt.Errorf("output handling failed: %w", err)
	}

	// Display file mappings
	fmt.Println("\nInput Path:", args.InputPath.Path)
	for _, f := range args.InputPath.Files {
		fmt.Printf("\n- %s", f)
	}
	fmt.Println("\nOutput Path:", args.OutputPath.Path)
	for _, f := range args.OutputPath.Files {
		fmt.Printf("\n- %s", f)
	}

	// Check file count consistency
	if len(args.InputPath.Files) != len(args.OutputPath.Files) {
		return fmt.Errorf("input/output file count mismatch")
	}

	// Process each input file and generate JSON
	for i := range args.InputPath.Files {
		fin := args.InputPath.Files[i]
		fout := args.OutputPath.Files[i]

		if err := fileIO.MakeFile(fout); err != nil {
			return fmt.Errorf("failed to create output file %q: %w", fout, err)
		}

		fmt.Println("\nProcessing file:", fin)

		result, err := xmlTools.ParsePubmedArticleSet(fin)
		if err != nil {
			return fmt.Errorf("failed to parse XML %q: %w", fin, err)
		}

		if err := jsonTools.ConvertToJson(result, fout); err != nil {
			return fmt.Errorf("failed to convert to JSON %q: %w", fout, err)
		}
	}

	return nil
}
