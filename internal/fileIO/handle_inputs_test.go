package fileIO_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ashahide/pubparse/internal/fileIO"
)

func TestHandleInputs_SingleFile(t *testing.T) {
	tmpDir := t.TempDir()
	xmlPath := filepath.Join(tmpDir, "sample.xml")

	// create a fake XML file
	if err := os.WriteFile(xmlPath, []byte("<root></root>"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{Path: xmlPath},
	}

	if err := fileIO.HandleInputs(args); err != nil {
		t.Errorf("HandleInputs failed for single file: %v", err)
	}

	if len(args.InputPath.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(args.InputPath.Files))
	}
}

func TestHandleInputs_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "a.xml")
	file2 := filepath.Join(tmpDir, "b.xml")
	os.WriteFile(file1, []byte("<a></a>"), 0644)
	os.WriteFile(file2, []byte("<b></b>"), 0644)

	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{Path: tmpDir},
	}

	if err := fileIO.HandleInputs(args); err != nil {
		t.Errorf("HandleInputs failed for directory: %v", err)
	}

	if len(args.InputPath.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(args.InputPath.Files))
	}
}

func TestHandleInputs_MissingPath(t *testing.T) {
	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{Path: ""},
	}

	err := fileIO.HandleInputs(args)
	if err == nil {
		t.Error("expected error for missing input path, got nil")
	}
}

func TestHandleInputs_WrongExtension(t *testing.T) {
	tmpDir := t.TempDir()
	txtFile := filepath.Join(tmpDir, "not_xml.txt")
	os.WriteFile(txtFile, []byte("not an xml"), 0644)

	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{Path: tmpDir},
	}

	err := fileIO.HandleInputs(args)
	if err != nil {
		// expected to skip the file, but not fail unless no valid XML files
		t.Logf("Handled wrong extension correctly: %v", err)
	} else if len(args.InputPath.Files) > 0 {
		t.Errorf("expected no files, got %d", len(args.InputPath.Files))
	}
}

func TestHandleInputs_InvalidPath(t *testing.T) {
	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{Path: "/this/does/not/exist"},
	}

	err := fileIO.HandleInputs(args)
	if err == nil {
		t.Error("expected error for non-existent path, got nil")
	}
}
