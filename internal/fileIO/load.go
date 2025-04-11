package fileIO

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ashahide/pubparse/internal/data"
)

func LoadFilesInDir(dirInfo data.PathInfo) (data.PathInfo, error) {
	if !dirInfo.Info.IsDir() {
		// Treat single file as the only entry
		dirInfo.Files = append(dirInfo.Files, data.PathInfo{
			Path: dirInfo.Path,
			Info: dirInfo.Info,
		})
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
		fileInfo, err := VerifyInputPath(fullPath)
		if err != nil {
			return dirInfo, fmt.Errorf("failed to verify path for file %q: %w", fullPath, err)
		}

		dirInfo.Files = append(dirInfo.Files, data.PathInfo{
			Path: fullPath,
			Info: fileInfo,
		})
	}

	return dirInfo, nil
}
