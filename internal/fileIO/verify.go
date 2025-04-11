package fileIO

import (
	"fmt"
	"os"
	"strings"
)

func VerifyPath(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)

	// If the error is because it doesn't exist
	if os.IsNotExist(err) {
		// Only return error if it's a directory that doesn't exist
		if strings.HasSuffix(path, string(os.PathSeparator)) {
			return nil, fmt.Errorf("directory does not exist: %s", path)
		}
		// For files that don’t exist, allow it (return nil, no error)
		return nil, nil
	}

	if os.IsPermission(err) {
		return nil, fmt.Errorf("permission denied: %s", path)
	}

	if err != nil {
		return nil, fmt.Errorf("error checking path %s: %w", path, err)
	}

	return info, nil
}

func VerifyInputPath(path string) (os.FileInfo, error) {
	return VerifyPath(path)
}

func VerifyOutputPath(path string) (os.FileInfo, error) {
	info, err := VerifyPath(path)
	if err != nil {
		return nil, err
	}

	if info == nil {
		// File doesn't exist — try creating it
		f, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("could not create output file %q: %w", path, err)
		}
		defer f.Close()

		// Re-stat after creation
		info, err = os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("file created but stat failed: %w", err)
		}
	}

	if info.IsDir() {
		return nil, fmt.Errorf("output path must be a file, got directory: %s", path)
	}

	return info, nil
}
