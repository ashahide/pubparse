package fileIO

import (
	"fmt"
	"os"
	"path/filepath"
)

// HandleOutputs prepares the output directory and file paths for writing JSON files.
//
// It performs the following steps:
//  1. Determines the appropriate "process" output directory.
//  2. Ensures the output directory exists (or creates it).
//  3. Generates one-to-one output .json file paths corresponding to the input files.
//  4. Verifies write access by attempting to create each file.
//  5. Verifies that the output directory exists and is accessible.
//
// The resulting output paths are stored in args.OutputPath.
func HandleOutputs(args *Arguments) error {
	// Step 1: Get the output directory (sibling "process" directory to input)
	processDir, err := GetProcessDir(args.InputPath.Path, args.InputPath.Info)
	if err != nil {
		return err
	}

	// Step 2: Ensure that the output directory exists
	if err := EnsureDir(processDir); err != nil {
		return err
	}

	// Step 3: Create full paths for each output .json file
	outputFiles, err := GenerateJSONFilePaths(args.InputPath.Files, processDir)
	if err != nil {
		return err
	}

	// Step 4: Ensure we can create/write each output file
	if err := VerifyWriteAccess(outputFiles); err != nil {
		return err
	}

	// Step 5: Get info about the output directory itself
	info, err := VerifyPath(processDir, "")
	if err != nil {
		return fmt.Errorf("failed to verify output path %q: %w", processDir, err)
	}

	// Store output path info in args
	args.OutputPath = PathInfo{
		Path:  processDir,
		Files: outputFiles,
		Info:  info,
	}

	return nil
}

// GetProcessDir determines the parallel "process" output directory for a given input path.
//
// If the input is a directory, it returns a sibling "process" directory.
// If the input is a file, it returns a "process" directory that is a sibling of the file’s parent.
//
// Example:
//
//	Input: /data/input/article.xml
//	Output: /data/process
func GetProcessDir(inputPath string, inputInfo os.FileInfo) (string, error) {
	if inputInfo.IsDir() {
		// For input directory, create sibling 'process' directory
		return filepath.Join(filepath.Dir(inputPath), "process"), nil
	}

	// For input file, create sibling to the parent of the file
	parentDir := filepath.Dir(inputPath)
	grandParent := filepath.Dir(parentDir)
	return filepath.Join(grandParent, "process"), nil
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
