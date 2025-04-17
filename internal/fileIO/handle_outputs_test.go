package fileIO_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ashahide/pubparse/internal/fileIO"
)

// --- getProcessDir ---

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

func TestGetProcessDir_File(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "file.xml")
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

// --- ensureDir ---

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

// --- verifyWriteAccess ---

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

// --- HandleOutputs ---

func TestHandleOutputs_Success(t *testing.T) {
	// Set up input
	tmpDir := t.TempDir()
	xmlFile := filepath.Join(tmpDir, "article.xml")
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

	if len(args.OutputPath.Files) != 1 {
		t.Errorf("expected 1 output file, got %d", len(args.OutputPath.Files))
	}

	if _, err := os.Stat(args.OutputPath.Files[0]); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestHandleOutputs_InvalidOutputPath(t *testing.T) {
	// Use an obviously bad path that shouldn't be writable
	badPath := "/this/should/not/exist/ever/file.xml"

	args := &fileIO.Arguments{
		InputPath: fileIO.PathInfo{
			Path:  badPath,
			Files: []string{badPath},
			Info:  &fakeFileInfo{name: "file.xml", isDir: false},
		},
	}

	err := fileIO.HandleOutputs(args)
	if err == nil {
		t.Error("expected HandleOutputs to fail on bad path, got nil")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}

// --- helper ---

type fakeFileInfo struct {
	name  string
	isDir bool
}

func (f *fakeFileInfo) Name() string           { return f.name }
func (f *fakeFileInfo) Size() int64            { return 0 }
func (f *fakeFileInfo) Mode() os.FileMode      { return 0644 }
func (f *fakeFileInfo) ModTime() (t time.Time) { return }
func (f *fakeFileInfo) IsDir() bool            { return f.isDir }
func (f *fakeFileInfo) Sys() interface{}       { return nil }
