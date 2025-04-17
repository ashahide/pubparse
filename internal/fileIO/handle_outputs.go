package fileIO

import (
	"fmt"
	"os"
	"path/filepath"
)

func HandleOutputs(args *Arguments) error {
	processDir, err := GetProcessDir(args.InputPath.Path, args.InputPath.Info)
	if err != nil {
		return err
	}

	if err := EnsureDir(processDir); err != nil {
		return err
	}

	outputFiles, err := GenerateJSONFilePaths(args.InputPath.Files, processDir)
	if err != nil {
		return err
	}

	if err := VerifyWriteAccess(outputFiles); err != nil {
		return err
	}

	info, err := VerifyPath(processDir, "")
	if err != nil {
		return fmt.Errorf("failed to verify output path %q: %w", processDir, err)
	}

	args.OutputPath = PathInfo{
		Path:  processDir,
		Files: outputFiles,
		Info:  info,
	}

	return nil
}

func GetProcessDir(inputPath string, inputInfo os.FileInfo) (string, error) {
	if inputInfo.IsDir() {
		return filepath.Join(filepath.Dir(inputPath), "process"), nil
	}
	parentDir := filepath.Dir(inputPath)
	grandParent := filepath.Dir(parentDir)
	return filepath.Join(grandParent, "process"), nil
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func VerifyWriteAccess(paths []string) error {
	for _, p := range paths {
		dir := filepath.Dir(p)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create subdirectory for %q: %w", p, err)
		}
		f, err := os.Create(p)
		if err != nil {
			return fmt.Errorf("cannot create output file %q: %w", p, err)
		}
		f.Close()
	}
	return nil
}
