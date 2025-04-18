package fileIO_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ashahide/pubparse/internal/fileIO"
)

//
// ---------------------- GetOutputDir ----------------------
//

// TestGetOutputDir_Directory verifies that GetOutputDir returns the correct
// auto-generated output directory when the input path is a directory.
// The expected format is: <input_parent>/processed_<dirname>
func TestGetOutputDir_Directory(t *testing.T) {
	tmp := t.TempDir()

	info, err := os.Stat(tmp)
	if err != nil {
		t.Fatalf("failed to stat temp dir: %v", err)
	}

	result, err := fileIO.GetOutputDir(tmp, info, "", nil)
	if err != nil {
		t.Errorf("GetOutputDir returned error: %v", err)
	}

	expected := filepath.Join(filepath.Dir(tmp), "processed_"+filepath.Base(tmp))
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// TestGetOutputDir_File verifies that GetOutputDir returns the correct
// auto-generated output directory when the input path is a file.
// The expected format is: <file_grandparent>/processed_<parent_dir>
func TestGetOutputDir_File(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "file.xml")

	// Create a dummy file to test against
	if err := os.WriteFile(tmpFile, []byte("hi"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("stat file failed: %v", err)
	}

	parentDir := filepath.Base(filepath.Dir(tmpFile))
	expected := filepath.Join(filepath.Dir(filepath.Dir(tmpFile)), "processed_"+parentDir)

	result, err := fileIO.GetOutputDir(tmpFile, info, "", nil)
	if err != nil {
		t.Errorf("GetOutputDir returned error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

//
// ---------------------- EnsureDir ----------------------
//

// TestEnsureDir_CreatesDirectory ensures that EnsureDir creates
// nested directory structures successfully.
func TestEnsureDir_CreatesDirectory(t *testing.T) {
	tmp := t.TempDir()
	newPath := filepath.Join(tmp, "subdir", "more")

	err := fileIO.EnsureDir(newPath)
	if err != nil {
		t.Errorf("EnsureDir failed: %v", err)
	}

	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Errorf("directory not created: %s", newPath)
	}
}

//
// ---------------------- VerifyWriteAccess ----------------------
//

// TestVerifyWriteAccess_CreatesFiles verifies that VerifyWriteAccess
// creates each file successfully and ensures write permissions.
func TestVerifyWriteAccess_CreatesFiles(t *testing.T) {
	tmp := t.TempDir()
	files := []string{
		filepath.Join(tmp, "out1.json"),
		filepath.Join(tmp, "out2.json"),
	}

	err := fileIO.VerifyWriteAccess(files)
	if err != nil {
		t.Errorf("VerifyWriteAccess failed: %v", err)
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

// TestHandleOutputs_Success ensures the full HandleOutputs pipeline works as expected.
// It should:
// - Auto-generate an output directory
// - Create one .json file path per input
// - Validate output path naming convention
// - Actually create the output files
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
		// Leave OutputPath.Path blank to trigger auto-generation
	}

	if err := fileIO.HandleOutputs(args); err != nil {
		t.Errorf("HandleOutputs failed: %v", err)
	}

	// Expect one .json file
	if len(args.OutputPath.Files) != 1 {
		t.Errorf("expected 1 output file, got %d", len(args.OutputPath.Files))
	}

	// Confirm auto-naming follows processed_ pattern
	if !strings.HasPrefix(filepath.Base(args.OutputPath.Path), "processed_") {
		t.Errorf("expected output path to start with 'processed_', got %q", args.OutputPath.Path)
	}

	// Check that file was created
	if _, err := os.Stat(args.OutputPath.Files[0]); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

// TestHandleOutputs_InvalidOutputPath simulates a failure scenario where the
// input path is invalid or inaccessible. It should return an error.
func TestHandleOutputs_InvalidOutputPath(t *testing.T) {
	// Simulated invalid path
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
