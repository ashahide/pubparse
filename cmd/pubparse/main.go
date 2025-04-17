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

	// Define command line flags
	flag.StringVar(&args.InputPath.Path, "i", "", "Path to the input file or directory of files")
	flag.Parse()

	// Process Inputs
	fileIO.HandleInputs(&args)

	// Process Outputs
	fileIO.HandleOutputs(&args)

	// Print file info
	fmt.Println("\nInput Path:", args.InputPath.Path)
	for _, f := range args.InputPath.Files {
		fmt.Printf("\n- %s", f)
	}

	fmt.Println("\nOutput Path:", args.OutputPath.Path)
	for _, f := range args.OutputPath.Files {
		fmt.Printf("\n- %s", f)
	}

	if len(args.InputPath.Files) != len(args.OutputPath.Files) {
		return fmt.Errorf("input/output file count mismatch")
	}

	for i := range args.InputPath.Files {
		fin := args.InputPath.Files[i]
		fout := args.OutputPath.Files[i]

		// Make output files
		if err := fileIO.MakeFile(fout); err != nil {
			return fmt.Errorf("failed to create output file %q: %w", fout, err)
		}

		fmt.Println("\nProcessing file:", fin)

		result, err := xmlTools.ParsePubmedArticleSet(fin)
		if err != nil {
			log.Fatal(err)
		}

		err = jsonTools.ConvertToJson(result, fout)
		if err != nil {
			log.Fatal(err)
		}

	}

	return nil
}
