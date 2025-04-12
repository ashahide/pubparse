package fileIO

import (
	"fmt"
	"os"
)

func VerifyPath(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	switch {
	case os.IsNotExist(err):
		return nil, fmt.Errorf("path does not exist: %s", path)
	case os.IsPermission(err):
		return nil, fmt.Errorf("permission denied: %s", path)
	case err != nil:
		return nil, fmt.Errorf("error checking path %s: %w", path, err)
	default:
		return info, nil
	}
}
