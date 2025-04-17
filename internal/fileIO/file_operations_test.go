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

// TestChangeExtension validates the ChangeExtension function,
// which replaces or adds a file extension to a given path.
func TestChangeExtension(t *testing.T) {
	// Define a set of test cases using table-driven style.
	tests := []struct {
		name     string // Name of the test case for readability in test output
		path     string // Original file path
		newExt   string // New extension to apply
		expected string // Expected result after extension change
	}{
		{"xml to json", "test.xml", ".json", "test.json"},         // Typical extension replacement
		{"xml to json - no dot", "test.xml", "json", "test.json"}, // Handles missing dot in extension
		{"empty path", "", ".json", ".json"},                      // Edge case: empty input path
	}

	// Loop through each test case
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call the function under test
			result := fileIO.ChangeExtension(test.path, test.newExt)

			// Compare the actual result with the expected result
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}
