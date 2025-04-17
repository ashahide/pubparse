package fileIO

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func GenerateJSONFilePaths(inputFiles []string, outputDir string) ([]string, error) {
	var outputPaths []string
	for _, inputFile := range inputFiles {
		base := filepath.Base(inputFile)
		jsonFile := ChangeExtension(base, "json")
		outputPath := filepath.Join(outputDir, jsonFile)
		outputPaths = append(outputPaths, outputPath)
	}
	return outputPaths, nil
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
