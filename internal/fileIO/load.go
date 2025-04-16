package fileIO

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ashahide/pubparse/internal/customErrors"
)

func LoadFilesInDir(dirInfo PathInfo, desiredTypeExt string) (PathInfo, error) {
	if !dirInfo.Info.IsDir() {
		// Treat single file as the only entry
		entry, err := os.Stat(dirInfo.Path)
		if err != nil {
			return dirInfo, fmt.Errorf("could not read file %q: %w", dirInfo.Path, err)
		}
		dirInfo.Files = append(dirInfo.Files, entry)
		return dirInfo, nil
	}

	entries, err := os.ReadDir(dirInfo.Path)
	if err != nil {
		return dirInfo, fmt.Errorf("could not read directory %q: %w", dirInfo.Path, err)
	}

	if len(entries) == 0 {
		return dirInfo, fmt.Errorf("no files found in directory: %s", dirInfo.Path)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirInfo.Path, entry.Name())
		fileInfo, err := VerifyPath(fullPath, desiredTypeExt)
		if err != nil {
			// Check if it's a WrongExtensionError and skip
			var extErr *customErrors.WrongExtensionError
			if errors.As(err, &extErr) {
				continue
			}
			// For any other error, return immediately
			return dirInfo, fmt.Errorf("failed to verify path for file %q: %w", fullPath, err)
		}

		dirInfo.Files = append(dirInfo.Files, fileInfo)
	}

	return dirInfo, nil
}
