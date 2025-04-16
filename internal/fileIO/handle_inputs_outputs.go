package fileIO

import (
	"fmt"
	"os"
	"path/filepath"
)

func HandleInputs(args *Arguments) error {
	var err error

	// Validate flags
	if args.InputPath.Path == "" {
		return fmt.Errorf("input path is required")
	}

	// Resolve absolute path
	args.InputPath.Path, err = filepath.Abs(args.InputPath.Path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %q: %w", args.InputPath.Path, err)
	}

	// Verify input path
	args.InputPath.Info, err = VerifyPath(args.InputPath.Path, "")
	if err != nil {
		return fmt.Errorf("failed to verify input path %q: %w", args.InputPath.Path, err)
	}

	// Load files from input directory
	args.InputPath, err = LoadFilesInDir(args.InputPath, "xml")
	if err != nil {
		return fmt.Errorf("failed to load input files from %q: %w", args.InputPath.Path, err)
	}

	// Just use the input dir
	args.OutputPath = PathInfo{
		Path:  args.InputPath.Path,
		Info:  args.InputPath.Info,
		Files: append([]os.FileInfo(nil), args.InputPath.Files...), // safe copy
	}

	// Update the path to be a json
	if !args.InputPath.Info.IsDir() {
		args.OutputPath.Path = ChangeExtension(args.OutputPath.Path, "json")
	}

	// Convert any input xml paths to output json paths
	args.OutputPath.Files, err = ConvertXMLtoJSON(args.OutputPath.Files, args.OutputPath.Path)
	if err != nil {
		return fmt.Errorf("failed to convert input file xml paths to output file jsons %q: %w", args.InputPath.Files, err)
	}

	// Verify output paths
	if args.OutputPath.Info, err = VerifyPath(args.OutputPath.Path, ""); err != nil {
		return fmt.Errorf("failed to verify output path %q: %w", args.OutputPath.Path, err)
	}

	return nil
}
