package fileIO_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ashahide/pubparse/internal/fileIO"
)

// TestHandleInputs_SingleFile verifies that HandleInputs correctly loads a single .xml file
// provided directly as the input path.
func TestHandleInputs_SingleFile(t *testing.T) {
	tmpDir := t.TempDir()
	xmlPath := filepath.Join(tmpDir, "sample.xml")

	// Create a dummy XML file
	if err := os.WriteFile(xmlPath, []byte("<root></root>"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{Path: xmlPath},
	}

	// Run input handler
	if err := fileIO.HandleInputs(args); err != nil {
		t.Errorf("HandleInputs failed for single file: %v", err)
	}

	// Ensure exactly one file was detected
	if len(args.InputPath.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(args.InputPath.Files))
	}
}

// TestHandleInputs_Directory verifies that HandleInputs loads multiple .xml files
// when a directory is provided as the input path.
func TestHandleInputs_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two dummy XML files in the directory
	file1 := filepath.Join(tmpDir, "a.xml")
	file2 := filepath.Join(tmpDir, "b.xml")
	os.WriteFile(file1, []byte("<a></a>"), 0644)
	os.WriteFile(file2, []byte("<b></b>"), 0644)

	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{Path: tmpDir},
	}

	// Run input handler
	if err := fileIO.HandleInputs(args); err != nil {
		t.Errorf("HandleInputs failed for directory: %v", err)
	}

	// Expect both files to be loaded
	if len(args.InputPath.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(args.InputPath.Files))
	}
}

// TestHandleInputs_MissingPath ensures that HandleInputs fails if no input path is provided.
func TestHandleInputs_MissingPath(t *testing.T) {
	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{Path: ""},
	}

	// Expect an error due to missing path
	err := fileIO.HandleInputs(args)
	if err == nil {
		t.Error("expected error for missing input path, got nil")
	}
}

// TestHandleInputs_WrongExtension verifies that files with unsupported extensions
// are skipped, and no input files are loaded.
func TestHandleInputs_WrongExtension(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a non-XML file in the directory
	txtFile := filepath.Join(tmpDir, "not_xml.txt")
	os.WriteFile(txtFile, []byte("not an xml"), 0644)

	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{Path: tmpDir},
	}

	// Run input handler
	err := fileIO.HandleInputs(args)

	// It may still succeed if logic allows no .xml files, but shouldn't return any files
	if err != nil {
		t.Logf("Handled wrong extension correctly: %v", err)
	} else if len(args.InputPath.Files) > 0 {
		t.Errorf("expected no files, got %d", len(args.InputPath.Files))
	}
}

// TestHandleInputs_InvalidPath ensures that HandleInputs returns an error
// when the input path does not exist.
func TestHandleInputs_InvalidPath(t *testing.T) {
	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{Path: "/this/does/not/exist"},
	}

	// Run input handler on an invalid path
	err := fileIO.HandleInputs(args)
	if err == nil {
		t.Error("expected error for non-existent path, got nil")
	}
}
