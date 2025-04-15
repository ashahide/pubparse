package fileIO_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ashahide/pubparse/internal/fileIO"
)

// TestMakeFile validates that MakeFile correctly creates or recreates files.
func TestMakeFile(t *testing.T) {
	// Create a temporary directory for safe, isolated test files
	tmpDir := t.TempDir()

	t.Run("create new file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "newfile.txt")

		// Call the function
		err := fileIO.MakeFile(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check if the file now exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("file was not created at: %s", path)
		}
	})

	t.Run("recreate existing file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "existing.txt")

		// Create file manually
		err := os.WriteFile(path, []byte("original"), 0644)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		// Re-create it via MakeFile
		err = fileIO.MakeFile(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Confirm the file exists
		_, err = os.Stat(path)
		if err != nil {
			t.Errorf("file disappeared: %v", err)
		}
	})
}
