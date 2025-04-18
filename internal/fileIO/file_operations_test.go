package fileIO_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ashahide/pubparse/internal/fileIO"
)

//
// ------------------------ Test: MakeFile ------------------------
//

// TestMakeFile verifies that MakeFile correctly creates new files,
// and also replaces existing files without raising errors.
func TestMakeFile(t *testing.T) {
	// Create a temporary directory for isolated test artifacts
	tmpDir := t.TempDir()

	t.Run("create new file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "newfile.txt")

		// Attempt to create the new file
		err := fileIO.MakeFile(path)
		if err != nil {
			t.Fatalf("unexpected error while creating file: %v", err)
		}

		// Verify the file now exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("file was not created at: %s", path)
		}
	})

	t.Run("recreate existing file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "existing.txt")

		// Manually create the file
		err := os.WriteFile(path, []byte("original content"), 0644)
		if err != nil {
			t.Fatalf("setup failed: could not create initial file: %v", err)
		}

		// Call MakeFile to overwrite it
		err = fileIO.MakeFile(path)
		if err != nil {
			t.Fatalf("unexpected error while recreating file: %v", err)
		}

		// Confirm the file still exists (recreated)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("file was unexpectedly missing: %v", err)
		}
	})
}

//
// ------------------------ Test: ChangeExtension ------------------------
//

// TestChangeExtension ensures that ChangeExtension correctly changes or adds
// file extensions to input paths.
//
// It handles paths with existing extensions, no extensions, or missing dots in new extensions.
func TestChangeExtension(t *testing.T) {
	// Table of test cases
	tests := []struct {
		name     string // Descriptive name for subtest
		path     string // Input file path
		newExt   string // New extension to apply
		expected string // Expected result
	}{
		{
			name:     "xml to json",
			path:     "test.xml",
			newExt:   ".json",
			expected: "test.json",
		},
		{
			name:     "xml to json - no dot",
			path:     "test.xml",
			newExt:   "json",
			expected: "test.json",
		},
		{
			name:     "empty path",
			path:     "",
			newExt:   ".json",
			expected: ".json",
		},
	}

	// Run each test case in isolation
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Execute the function under test
			result := fileIO.ChangeExtension(test.path, test.newExt)

			// Validate the result
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}
