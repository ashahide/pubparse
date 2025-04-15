package fileIO

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ashahide/pubparse/internal/customErrors"
)

// VerifyPath checks if a given file or directory path is valid and optionally verifies its extension.
//
// Arguments:
//   - path: the file or directory path to verify.
//   - fileType: an optional file extension (e.g., "json", ".txt") to validate against the path.
//
// Returns:
//   - os.FileInfo describing the file or directory if it exists and passes validation.
//   - error if the path does not exist, is not accessible, or fails the fileType check.
func VerifyPath(path string, fileType string) (os.FileInfo, error) {
	// Convert the input path to an absolute path for consistent evaluation
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("error converting to absolute path %s: %w", path, err)
	}

	// Check if the path exists and gather file metadata
	info, err := os.Stat(absPath)
	switch {
	case os.IsNotExist(err):
		// File or directory does not exist
		return nil, fmt.Errorf("path does not exist: %s", absPath)
	case os.IsPermission(err):
		// Lack of permissions to access the file or directory
		return nil, fmt.Errorf("permission denied: %s", absPath)
	case err != nil:
		// Any other error while accessing the path
		return nil, fmt.Errorf("error checking path %s: %w", absPath, err)
	}

	// If a fileType was specified, verify the path has the correct extension
	if fileType != "" {
		// Normalize fileType (remove leading dot and lowercase it)
		expectedExt := "." + strings.TrimPrefix(strings.ToLower(fileType), ".")

		// Extract and normalize the actual file extension
		actualExt := strings.ToLower(filepath.Ext(absPath))

		// Compare the expected and actual extensions
		if actualExt != expectedExt {
			return nil, &customErrors.WrongExtensionError{
				Expected: expectedExt,
				Actual:   actualExt,
			}
		}
	}

	// Return the FileInfo if all checks passed
	return info, nil
}
