package fileIO

import (
	"fmt"
	"os"
	"path/filepath"
)

// HandleInputs orchestrates validation, file loading, and output preparation.
// It updates the Arguments struct with validated paths and converted output filenames.
func HandleInputs(args *Arguments) error {
	var err error

	// Step 1: Validate the input path and store absolute version in args
	if err = validateInputPath(args); err != nil {
		return err
	}

	// Step 2: Load XML files from the validated input path
	if err = populateInputFiles(args); err != nil {
		return err
	}

	// Step 3: Prepare output paths and filenames based on the inputs
	if err = prepareOutputPath(args); err != nil {
		return err
	}

	return nil
}

// validateInputPath resolves and verifies the input path.
// It ensures the path exists and is accessible.
func validateInputPath(args *Arguments) error {
	if args.InputPath.Path == "" {
		return fmt.Errorf("input path is required")
	}

	// Resolve the input path to an absolute path
	absPath, err := filepath.Abs(args.InputPath.Path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %q: %w", args.InputPath.Path, err)
	}
	args.InputPath.Path = absPath

	// Verify that the input path is accessible
	args.InputPath.Info, err = VerifyPath(absPath, "")
	if err != nil {
		return fmt.Errorf("failed to verify input path %q: %w", absPath, err)
	}

	return nil
}

// populateInputFiles loads the XML files from the input directory.
// If the input path is a single file, it wraps it as a list.
func populateInputFiles(args *Arguments) error {
	var err error

	// Attempt to load all XML files in the directory (or single file)
	args.InputPath, err = LoadFilesInDir(args.InputPath, "xml")
	if err != nil {
		return fmt.Errorf("failed to load input files from %q: %w", args.InputPath.Path, err)
	}
	return nil
}

// prepareOutputPath creates output file paths by changing the extension of each input file to .json.
// It sets the output path and verifies the output destination.
func prepareOutputPath(args *Arguments) error {
	// Use input directory as base for outputs
	args.OutputPath = PathInfo{
		Path:  args.InputPath.Path,
		Info:  args.InputPath.Info,
		Files: append([]os.FileInfo(nil), args.InputPath.Files...),
	}

	if !args.InputPath.Info.IsDir() {
		// Use the parent directory, not the file name with .json
		args.OutputPath.Path = filepath.Dir(args.OutputPath.Path)
	}

	var err error
	args.OutputPath.Files, err = ConvertXMLtoJSON(args.OutputPath.Files, args.OutputPath.Path)
	if err != nil {
		return fmt.Errorf("failed to convert input file names to JSON: %w", err)
	}

	args.OutputPath.Info, err = VerifyPath(args.OutputPath.Path, "")
	if err != nil {
		return fmt.Errorf("failed to verify output path %q: %w", args.OutputPath.Path, err)
	}

	return nil
}
