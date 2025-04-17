package fileIO

import (
	"fmt"
	"path/filepath"
)

// HandleInputs is the main entry point for processing input file arguments.
//
// It performs the following:
//  1. Resolves and validates the user-provided input path.
//  2. Loads all valid `.xml` files from the input path (either a directory or single file).
//
// On success, it populates the InputPath field in `args` with:
//   - Absolute path
//   - Path metadata (os.FileInfo)
//   - Slice of matching XML file paths
func HandleInputs(args *Arguments) error {
	var err error

	// Step 1: Ensure the input path is valid and resolved
	if err = validateInputPath(args); err != nil {
		return err
	}

	// Step 2: Load all .xml files from the input path
	if err = populateInputFiles(args); err != nil {
		return err
	}

	return nil
}

// validateInputPath resolves the provided input path to an absolute path,
// and verifies that it exists and is accessible.
//
// If validation is successful, it updates:
//   - args.InputPath.Path with the absolute path
//   - args.InputPath.Info with the file/directory metadata
func validateInputPath(args *Arguments) error {
	if args.InputPath.Path == "" {
		return fmt.Errorf("input path is required")
	}

	// Resolve relative or symbolic paths to absolute
	absPath, err := filepath.Abs(args.InputPath.Path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %q: %w", args.InputPath.Path, err)
	}
	args.InputPath.Path = absPath

	// Validate existence and accessibility of the resolved path
	args.InputPath.Info, err = VerifyPath(absPath, "")
	if err != nil {
		return fmt.Errorf("failed to verify input path %q: %w", absPath, err)
	}

	return nil
}

// populateInputFiles loads the list of `.xml` files from the provided input path.
//
// - If the input is a single file, it wraps it in a list.
// - If the input is a directory, it filters and loads all `.xml` files.
// On success, it updates args.InputPath.Files.
func populateInputFiles(args *Arguments) error {
	var err error

	// Load valid XML files (or the single file itself)
	args.InputPath, err = LoadFilesInDir(args.InputPath, "xml")
	if err != nil {
		return fmt.Errorf("failed to load input files from %q: %w", args.InputPath.Path, err)
	}

	return nil
}
