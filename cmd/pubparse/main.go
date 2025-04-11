package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ashahide/pubparse/internal/data"
	"github.com/ashahide/pubparse/internal/fileIO"
)

func main() {
	if err := run(); err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
}

func run() error {
	var args data.Arguments

	// Define command line flags
	flag.StringVar(&args.InputPath.Path, "i", "", "Path to the input file or directory of files")
	flag.StringVar(&args.OutputPath.Path, "o", "", "Path to the output file or directory of files")
	flag.Parse()

	// Validate flags
	if args.InputPath.Path == "" {
		return fmt.Errorf("input path is required")
	}
	if args.OutputPath.Path == "" {
		return fmt.Errorf("output path is required")
	}

	// Verify paths
	var err error
	if args.InputPath.Info, err = fileIO.VerifyInputPath(args.InputPath.Path); err != nil {
		return fmt.Errorf("failed to verify input path %q: %w", args.InputPath.Path, err)
	}
	if args.OutputPath.Info, err = fileIO.VerifyOutputPath(args.OutputPath.Path); err != nil {
		return fmt.Errorf("failed to verify output path %q: %w", args.OutputPath.Path, err)
	}

	// Load files
	if args.InputPath, err = fileIO.LoadFilesInDir(args.InputPath); err != nil {
		return fmt.Errorf("failed to load input files from %q: %w", args.InputPath.Path, err)
	}
	if args.OutputPath, err = fileIO.LoadFilesInDir(args.OutputPath); err != nil {
		return fmt.Errorf("failed to load output files from %q: %w", args.OutputPath.Path, err)
	}

	// Print file info
	fmt.Println("Input Path:", args.InputPath.Path)
	for _, f := range args.InputPath.Files {
		fmt.Printf(" - %s (%d bytes)\n", f.Path, f.Info.Size())
	}

	fmt.Println("Output Path:", args.OutputPath.Path)
	for _, f := range args.OutputPath.Files {
		fmt.Printf(" - %s (%d bytes)\n", f.Path, f.Info.Size())
	}

	return nil
}
