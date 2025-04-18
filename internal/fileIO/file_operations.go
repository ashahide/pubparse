package fileIO

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MakeFile creates a file at the specified path, overwriting it if it already exists.
//
// - If the file exists, it deletes and recreates it.
// - If the file does not exist, it simply creates it.
// - If any unexpected error occurs (e.g., permission denied), it returns an error.
func MakeFile(path string) error {
	_, err := os.Stat(path)

	if os.IsExist(err) {
		// File exists, remove and recreate it
		os.Remove(path)
		os.Create(path)
	} else if os.IsNotExist(err) {
		// File does not exist, just create it
		os.Create(path)
	} else if err != nil {
		// Any other error while checking file status
		return err
	}

	return nil
}

// GenerateJSONFilePaths takes a list of input XML file paths and returns a list
// of corresponding .json file paths in the given output directory.
//
// It converts each input file’s base name to a .json extension and appends it to outputDir.
func GenerateJSONFilePaths(inputFiles []string, outputDir string) ([]string, error) {
	var outputPaths []string

	for _, inputFile := range inputFiles {
		// Extract just the filename (e.g., "article.xml" → "article.json")
		base := filepath.Base(inputFile)
		jsonFile := ChangeExtension(base, "json")

		// Create full path in the output directory
		outputPath := filepath.Join(outputDir, jsonFile)
		outputPaths = append(outputPaths, outputPath)
	}

	return outputPaths, nil
}

// ChangeExtension replaces the file extension of a given path with a new one.
//
// If the path has no extension, the new extension is simply appended.
// The new extension can be passed with or without a leading dot.
func ChangeExtension(path, newExt string) string {
	ext := filepath.Ext(path)
	if ext != "" {
		// Remove old extension
		path = path[:len(path)-len(ext)]
	}

	// Ensure the new extension starts with "."
	if !strings.HasPrefix(newExt, ".") {
		newExt = "." + newExt
	}

	return path + newExt
}

// ConvertXMLtoJSON generates empty .json output files that match the given input XML files.
//
// - For each input os.FileInfo, it converts the filename to .json
// - It creates an empty file at the output location
// - It returns a slice of os.FileInfo for the created files
//
// This is mainly useful for testing or reserving output paths prior to actual JSON writing.
func ConvertXMLtoJSON(inputPath []os.FileInfo, outputDir string) ([]os.FileInfo, error) {
	var outputPath []os.FileInfo

	for _, entry := range inputPath {
		// Generate .json filename based on .xml file
		newFileName := ChangeExtension(entry.Name(), "json")
		fullPath := filepath.Join(outputDir, newFileName)

		// Create the output file
		if err := MakeFile(fullPath); err != nil {
			return nil, fmt.Errorf("failed to create output file %q: %w", fullPath, err)
		}

		// Get file info for the newly created file
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			return nil, err
		}

		outputPath = append(outputPath, fileInfo)
	}

	return outputPath, nil
}
