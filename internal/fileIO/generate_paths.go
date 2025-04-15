package fileIO

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ashahide/pubparse/internal/custom_errors"
	"github.com/ashahide/pubparse/internal/data"
)

func MakeFile(path string) error {

	_, err := os.Stat(path)
	if os.IsExist(err) {
		// File Exists
		os.Remove(path)
		os.Create(path)
	} else if os.IsNotExist(err) {
		// File Doesn't Exist
		os.Create(path)
	} else if err != nil {
		return err
	}

	return nil
}

// Change file extension
func ChangeExtension(path, newExt string) string {
	ext := filepath.Ext(path)
	if ext != "" {
		path = path[:len(path)-len(ext)]
	}
	if !strings.HasPrefix(newExt, ".") {
		newExt = "." + newExt
	}
	return path + newExt
}

func ConvertXMLtoJSON(input_path []os.FileInfo, outputDir string) ([]os.FileInfo, error) {
	var output_path []os.FileInfo

	for _, entry := range input_path {
		// Convert to .json
		newFileName := ChangeExtension(entry.Name(), "json")

		// Join with output directory
		fullPath := filepath.Join(outputDir, newFileName)

		fmt.Println("Creating file:", fullPath)

		if err := MakeFile(fullPath); err != nil {
			return nil, fmt.Errorf("failed to create output file %q: %w", fullPath, err)
		}

		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			return nil, err
		}

		output_path = append(output_path, fileInfo)
	}

	return output_path, nil
}

func LoadFilesInDir(dirInfo data.PathInfo, desiredTypeExt string) (data.PathInfo, error) {
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
			var extErr *custom_errors.WrongExtensionError
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

func VerifyPath(path string, fileType string) (os.FileInfo, error) {

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("error converting to absolute path %s: %w", path, err)
	}

	info, err := os.Stat(absPath)
	switch {
	case os.IsNotExist(err):
		return nil, fmt.Errorf("path does not exist: %s", absPath)
	case os.IsPermission(err):
		return nil, fmt.Errorf("permission denied: %s", absPath)
	case err != nil:
		return nil, fmt.Errorf("error checking path %s: %w", absPath, err)
	}

	if fileType != "" {
		expectedExt := "." + strings.TrimPrefix(strings.ToLower(fileType), ".")

		actualExt := strings.ToLower(filepath.Ext(absPath))
		if actualExt != expectedExt {
			return nil, &custom_errors.WrongExtensionError{
				Expected: expectedExt,
				Actual:   actualExt,
			}
		}
	}

	return info, nil
}
