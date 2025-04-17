package fileIO

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ashahide/pubparse/internal/customErrors"
)

// LoadFilesInDir inspects a given PathInfo and populates the .Files field
// with valid files that match the specified extension (e.g., "xml").
//
// Behavior:
//   - If the path is a single file, it wraps that file into the result.
//   - If it's a directory, it reads all files in that directory,
//     filters by extension, and verifies each one.
//   - If a file has the wrong extension, it is skipped without error.
//   - For any other error (e.g., permission, corrupt file), it returns immediately.
//
// Arguments:
//   - dirInfo: PathInfo containing the path and metadata (os.FileInfo).
//   - desiredTypeExt: Extension to match, e.g., "xml".
//
// Returns:
//   - Updated PathInfo with .Files populated with full file paths.
//   - Error if path is invalid, unreadable, or no matching files are found.
func LoadFilesInDir(dirInfo PathInfo, desiredTypeExt string) (PathInfo, error) {
	// If the input is a single file, treat it as a one-element list
	if !dirInfo.Info.IsDir() {
		_, err := os.Stat(dirInfo.Path)
		if err != nil {
			return dirInfo, fmt.Errorf("could not read file %q: %w", dirInfo.Path, err)
		}

		// Add the absolute path of the file
		dirInfo.Files = append(dirInfo.Files, dirInfo.Path)
		return dirInfo, nil
	}

	// Read all entries (files and dirs) in the directory
	entries, err := os.ReadDir(dirInfo.Path)
	if err != nil {
		return dirInfo, fmt.Errorf("could not read directory %q: %w", dirInfo.Path, err)
	}

	if len(entries) == 0 {
		return dirInfo, fmt.Errorf("no files found in directory: %s", dirInfo.Path)
	}

	for _, entry := range entries {
		// Construct full path to each entry
		fullPath := filepath.Join(dirInfo.Path, entry.Name())

		// Check if the path has the desired extension (e.g., ".xml")
		_, err := VerifyPath(fullPath, desiredTypeExt)
		if err != nil {
			// Skip files with wrong extensions without failing
			var extErr *customErrors.WrongExtensionError
			if errors.As(err, &extErr) {
				continue
			}
			// For any other verification error, exit early
			return dirInfo, fmt.Errorf("failed to verify path for file %q: %w", fullPath, err)
		}

		// Add the full path of valid file
		dirInfo.Files = append(dirInfo.Files, fullPath)
	}

	return dirInfo, nil
}
