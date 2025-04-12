package fileIO

import (
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

func ConvertXMLtoJSON(input_path []os.FileInfo) (output_path []os.FileInfo, err error) {

	for _, entry := range input_path {
		new_entry := ChangeExtension(entry.Name(), "json")

		MakeFile(new_entry)

		fileInfo, err := os.Stat(new_entry)

		if err != nil {
			return output_path, err
		}

		output_path = append(output_path, fileInfo)
	}

	return output_path, err
}
