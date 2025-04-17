package fileIO

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// HandleOutputs prepares the output directory and file paths for writing JSON files.
//
// It performs the following steps:
//  1. Determines the appropriate output directory (user-defined or auto-generated).
//  2. Ensures the output directory exists (or creates it).
//  3. Generates one-to-one output .json file paths corresponding to the input files.
//  4. Verifies write access by attempting to create each file.
//  5. Captures metadata about the output directory.
//
// The resulting output paths are stored in args.OutputPath.
func HandleOutputs(args *Arguments) error {
	var err error

	// Step 1: Generate output path if not provided
	if args.OutputPath.Path == "" {
		args.OutputPath.Path, err = GetOutputDir(
			args.InputPath.Path,
			args.InputPath.Info,
			"", // No output path provided
			nil,
		)
		if err != nil {
			return err
		}
	}

	// Step 2: Ensure that the output directory exists
	if err := EnsureDir(args.OutputPath.Path); err != nil {
		return err
	}

	// Step 3: Create full paths for each output .json file
	outputFiles, err := GenerateJSONFilePaths(args.InputPath.Files, args.OutputPath.Path)
	if err != nil {
		return err
	}

	// Step 4: Ensure we can create/write each output file
	if err := VerifyWriteAccess(outputFiles); err != nil {
		return err
	}

	// Step 5: Get info about the output directory itself
	info, err := VerifyPath(args.OutputPath.Path, "")
	if err != nil {
		return fmt.Errorf("failed to verify output path %q: %w", args.OutputPath.Path, err)
	}

	// Store output path info in args
	args.OutputPath.Info = info
	args.OutputPath.Files = outputFiles

	return nil
}

// GetOutputDir generates a default descriptive output directory path.
//
// If input is a directory: <input_parent>/processed_<input_dirname>
// If input is a file: <grandparent>/processed_<parent_dirname>
func GetOutputDir(inputPath string, inputInfo os.FileInfo, outputPath string, outputInfo os.FileInfo) (string, error) {
	var baseName string

	if inputInfo.IsDir() {
		baseName = filepath.Base(inputPath)
		return filepath.Join(filepath.Dir(inputPath), "processed_"+sanitizeName(baseName)), nil
	}

	parentDir := filepath.Dir(inputPath)
	baseName = filepath.Base(parentDir)
	grandParent := filepath.Dir(parentDir)
	return filepath.Join(grandParent, "processed_"+sanitizeName(baseName)), nil
}

// sanitizeName replaces spaces or separators with underscores to ensure safe directory names.
func sanitizeName(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}

// EnsureDir ensures that a directory exists, creating it and all parent directories if necessary.
//
// Returns an error if the directory cannot be created.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// VerifyWriteAccess attempts to create each of the given files to ensure that write access is permitted.
//
// For each file:
//   - Ensures the parent directory exists.
//   - Creates and immediately closes the file.
//   - If any operation fails, it returns an error.
//
// If all files can be written, it returns nil.
func VerifyWriteAccess(paths []string) error {
	for _, p := range paths {
		// Ensure parent directory exists
		dir := filepath.Dir(p)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create subdirectory for %q: %w", p, err)
		}

		// Try to create the file
		f, err := os.Create(p)
		if err != nil {
			return fmt.Errorf("cannot create output file %q: %w", p, err)
		}
		f.Close()
	}
	return nil
}
