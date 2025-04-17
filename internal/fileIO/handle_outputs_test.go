package fileIO_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ashahide/pubparse/internal/fileIO"
)

//
// ---------------------- getProcessDir ----------------------
//

// TestGetProcessDir_Directory verifies that getProcessDir correctly identifies
// the parent of a directory and returns a sibling "process" directory.
func TestGetProcessDir_Directory(t *testing.T) {
	tmp := t.TempDir()

	info, err := os.Stat(tmp)
	if err != nil {
		t.Fatalf("failed to stat temp dir: %v", err)
	}

	result, err := fileIO.GetProcessDir(tmp, info)
	if err != nil {
		t.Errorf("getProcessDir returned error: %v", err)
	}

	expected := filepath.Join(filepath.Dir(tmp), "process")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// TestGetProcessDir_File checks that getProcessDir returns the correct sibling
// "process" directory when given a file path as input.
func TestGetProcessDir_File(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "file.xml")

	// Create a dummy XML file
	if err := os.WriteFile(tmpFile, []byte("hi"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("stat file failed: %v", err)
	}

	result, err := fileIO.GetProcessDir(tmpFile, info)
	if err != nil {
		t.Errorf("getProcessDir returned error: %v", err)
	}

	expected := filepath.Join(filepath.Dir(filepath.Dir(tmpFile)), "process")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

//
// ---------------------- ensureDir ----------------------
//

// TestEnsureDir_CreatesDirectory ensures that EnsureDir creates nested directories as expected.
func TestEnsureDir_CreatesDirectory(t *testing.T) {
	tmp := t.TempDir()
	newPath := filepath.Join(tmp, "subdir", "more")

	err := fileIO.EnsureDir(newPath)
	if err != nil {
		t.Errorf("ensureDir failed: %v", err)
	}

	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Errorf("directory not created: %s", newPath)
	}
}

//
// ---------------------- verifyWriteAccess ----------------------
//

// TestVerifyWriteAccess_CreatesFiles checks that verifyWriteAccess can write to all given paths.
func TestVerifyWriteAccess_CreatesFiles(t *testing.T) {
	tmp := t.TempDir()
	files := []string{
		filepath.Join(tmp, "out1.json"),
		filepath.Join(tmp, "out2.json"),
	}

	err := fileIO.VerifyWriteAccess(files)
	if err != nil {
		t.Errorf("verifyWriteAccess failed: %v", err)
	}

	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected file not created: %s", f)
		}
	}
}

//
// ---------------------- HandleOutputs ----------------------
//

// TestHandleOutputs_Success verifies the full HandleOutputs pipeline for valid input,
// ensuring that process directory is created and JSON file paths are assigned.
func TestHandleOutputs_Success(t *testing.T) {
	tmpDir := t.TempDir()
	xmlFile := filepath.Join(tmpDir, "article.xml")

	// Create a dummy XML file
	_ = os.WriteFile(xmlFile, []byte("<xml></xml>"), 0644)

	info, _ := os.Stat(tmpDir)

	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{
			Path:  tmpDir,
			Info:  info,
			Files: []string{xmlFile},
		},
	}

	if err := fileIO.HandleOutputs(args); err != nil {
		t.Errorf("HandleOutputs failed: %v", err)
	}

	// Expect exactly one .json file output
	if len(args.OutputPath.Files) != 1 {
		t.Errorf("expected 1 output file, got %d", len(args.OutputPath.Files))
	}

	// Validate that the output file was created
	if _, err := os.Stat(args.OutputPath.Files[0]); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

// TestHandleOutputs_InvalidOutputPath simulates a failure when trying to write
// to an invalid path. It expects HandleOutputs to return an error.
func TestHandleOutputs_InvalidOutputPath(t *testing.T) {
	// Use a path that is extremely unlikely to exist
	badPath := "/this/should/not/exist/ever/file.xml"

	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{
			Path:  badPath,
			Files: []string{badPath},
			Info:  &fileIO.FakeFileInfo{NameVal: "file.xml", IsDirVal: false},
		},
	}

	err := fileIO.HandleOutputs(args)
	if err == nil {
		t.Error("expected HandleOutputs to fail on bad path, got nil")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}
