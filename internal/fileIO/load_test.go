package fileIO_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ashahide/pubparse/internal/fileIO"
)

//
// ------------------------ LoadFilesInDir Tests ------------------------
//

// TestLoadFilesInDir_SingleFile verifies that LoadFilesInDir correctly handles a case
// where the input path is a single .xml file. It should return the path wrapped in a slice.
func TestLoadFilesInDir_SingleFile(t *testing.T) {
	tmpDir := t.TempDir()
	xmlPath := filepath.Join(tmpDir, "single.xml")

	// Create a dummy XML file
	if err := os.WriteFile(xmlPath, []byte("<xml></xml>"), 0644); err != nil {
		t.Fatalf("failed to create XML file: %v", err)
	}

	// Get file info for the PathInfo struct
	info, _ := os.Stat(xmlPath)
	input := fileIO.PathInfo{
		Path: xmlPath,
		Info: info,
	}

	// Attempt to load the file
	result, err := fileIO.LoadFilesInDir(input, "xml")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(result.Files))
	}
}

// TestLoadFilesInDir_DirectoryAllValid verifies that LoadFilesInDir loads
// all valid .xml files from a given directory.
func TestLoadFilesInDir_DirectoryAllValid(t *testing.T) {
	tmpDir := t.TempDir()
	files := []string{"a.xml", "b.xml", "c.xml"}

	// Create valid XML files
	for _, name := range files {
		_ = os.WriteFile(filepath.Join(tmpDir, name), []byte("<xml/>"), 0644)
	}

	info, _ := os.Stat(tmpDir)
	input := fileIO.PathInfo{
		Path: tmpDir,
		Info: info,
	}

	// Attempt to load files
	result, err := fileIO.LoadFilesInDir(input, "xml")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Files) != len(files) {
		t.Errorf("expected %d files, got %d", len(files), len(result.Files))
	}
}

// TestLoadFilesInDir_MixedExtensions ensures that LoadFilesInDir correctly
// skips files that do not match the desired extension.
func TestLoadFilesInDir_MixedExtensions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create one valid and one invalid file
	os.WriteFile(filepath.Join(tmpDir, "valid.xml"), []byte("<ok/>"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "ignore.txt"), []byte("not xml"), 0644)

	info, _ := os.Stat(tmpDir)
	input := fileIO.PathInfo{
		Path: tmpDir,
		Info: info,
	}

	result, err := fileIO.LoadFilesInDir(input, "xml")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Files) != 1 {
		t.Errorf("expected 1 XML file, got %d", len(result.Files))
	}
}

// TestLoadFilesInDir_EmptyDir confirms that LoadFilesInDir returns an error
// when provided a directory that contains no files.
func TestLoadFilesInDir_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	info, _ := os.Stat(tmpDir)
	input := fileIO.PathInfo{
		Path: tmpDir,
		Info: info,
	}

	_, err := fileIO.LoadFilesInDir(input, "xml")
	if err == nil {
		t.Error("expected error on empty directory, got nil")
	}
}

// TestLoadFilesInDir_InvalidPath verifies that LoadFilesInDir fails
// when given a nonexistent or inaccessible path.
func TestLoadFilesInDir_InvalidPath(t *testing.T) {
	badPath := "/this/does/not/exist"

	input := fileIO.PathInfo{
		Path: badPath,
		Info: &fileIO.FakeFileInfo{NameVal: "bad", IsDirVal: true}, // simulate directory path
	}

	_, err := fileIO.LoadFilesInDir(input, "xml")
	if err == nil {
		t.Error("expected error on invalid path, got nil")
	}
}
