package fileIO

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ashahide/pubparse/internal/customErrors"
)

//
// ------------------------ VerifyPath ------------------------
//

// VerifyPath checks if the given file or directory path exists and is accessible,
// and optionally verifies that it has the expected file extension.
//
// Arguments:
//   - path: The file or directory path to validate (can be relative or absolute).
//   - fileType: An optional extension to check against (e.g., "json", ".xml").
//     If provided, the file must match this extension.
//
// Returns:
//   - os.FileInfo: Metadata describing the file or directory.
//   - error: If the path does not exist, is inaccessible, or has the wrong extension.
func VerifyPath(path string, fileType string) (os.FileInfo, error) {
	// Convert to absolute path for consistent filesystem access
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("error converting to absolute path %q: %w", path, err)
	}

	// Attempt to retrieve file or directory metadata
	info, err := os.Stat(absPath)
	switch {
	case os.IsNotExist(err):
		// File or directory does not exist
		return nil, fmt.Errorf("path does not exist: %s", absPath)
	case os.IsPermission(err):
		// Permission error while accessing the path
		return nil, fmt.Errorf("permission denied: %s", absPath)
	case err != nil:
		// Other filesystem errors
		return nil, fmt.Errorf("error checking path %s: %w", absPath, err)
	}

	// If an expected file type is provided, validate the extension
	if fileType != "" {
		// Normalize: ensure leading dot, and lowercase for consistent comparison
		expectedExt := "." + strings.TrimPrefix(strings.ToLower(fileType), ".")

		// Extract actual extension from the path
		actualExt := strings.ToLower(filepath.Ext(absPath))

		// If extensions don’t match, return a custom error
		if actualExt != expectedExt {
			return nil, &customErrors.WrongExtensionError{
				Expected: expectedExt,
				Actual:   actualExt,
			}
		}
	}

	// Path exists, is accessible, and matches expected extension (if given)
	return info, nil
}

//
// ------------------------ GenerateJSONFileInfos ------------------------
//

// GenerateJSONFileInfos creates a list of `.json` output file paths that correspond
// to a list of input files.
//
// For each input file path, the function:
//   - Extracts the base filename (without directory)
//   - Replaces its extension with `.json`
//   - Appends it to the specified output directory
//
// Example:
//
//	inputFiles: ["/data/a.xml", "/data/b.xml"]
//	outputDir:  "/results"
//	result:     ["/results/a.json", "/results/b.json"]
//
// Arguments:
//   - inputFiles: Slice of full paths to input XML files.
//   - outputDir:  Directory in which to place the corresponding .json files.
//
// Returns:
//   - Slice of .json file paths, each aligned with an input file.
//   - Error (always nil in current implementation; placeholder for future validation).
func GenerateJSONFileInfos(inputFiles []string, outputDir string) ([]string, error) {
	var outputPaths []string

	for _, entry := range inputFiles {
		// Extract base filename, change its extension to .json
		newFileName := ChangeExtension(entry, "json")

		// Combine with output directory to get full path
		fullPath := filepath.Join(outputDir, newFileName)

		// Add to result list
		outputPaths = append(outputPaths, fullPath)
	}

	return outputPaths, nil
}
