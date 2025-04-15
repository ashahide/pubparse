package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	var err error

	// Define command line flags
	flag.StringVar(&args.InputPath.Path, "i", "", "Path to the input file or directory of files")
	flag.Parse()

	// Validate flags
	if args.InputPath.Path == "" {
		return fmt.Errorf("input path is required")
	}

	args.InputPath.Path, err = filepath.Abs(args.InputPath.Path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %q: %w", args.InputPath.Path, err)
	}

	// Verify input paths
	if args.InputPath.Info, err = fileIO.VerifyPath(args.InputPath.Path, ""); err != nil {
		return fmt.Errorf("failed to verify input path %q: %w", args.InputPath.Path, err)
	}

	// Load files
	if args.InputPath, err = fileIO.LoadFilesInDir(args.InputPath, "xml"); err != nil {
		return fmt.Errorf("failed to load input files from %q: %w", args.InputPath.Path, err)
	}

	// Just use the input dir
	args.OutputPath = data.PathInfo{
		Path:  args.InputPath.Path,
		Info:  args.InputPath.Info,
		Files: append([]os.FileInfo(nil), args.InputPath.Files...), // safe copy
	}

	// Update the path to be a json
	if !args.InputPath.Info.IsDir() {
		args.OutputPath.Path = fileIO.ChangeExtension(args.OutputPath.Path, "json")
	}

	// Convert any input xml paths to output json paths
	args.OutputPath.Files, err = fileIO.ConvertXMLtoJSON(args.OutputPath.Files, args.OutputPath.Path)
	if err != nil {
		return fmt.Errorf("failed to convert input file xml paths to output file jsons %q: %w", args.InputPath.Files, err)
	}

	// Verify output paths
	if args.OutputPath.Info, err = fileIO.VerifyPath(args.OutputPath.Path, ""); err != nil {
		return fmt.Errorf("failed to verify output path %q: %w", args.OutputPath.Path, err)
	}

	// Print file info
	fmt.Println("Input Path:", args.InputPath.Path)
	for _, f := range args.InputPath.Files {
		fmt.Printf(" - %s (%d bytes)\n", f.Name(), f.Size())
	}

	fmt.Println("Output Path:", args.OutputPath.Path)
	for _, f := range args.OutputPath.Files {
		fmt.Printf(" - %s (%d bytes)\n", f.Name(), f.Size())
	}

	return nil
}
