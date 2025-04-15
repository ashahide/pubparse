package fileIO

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
