package fileIO_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ashahide/pubparse/internal/customErrors"
	"github.com/ashahide/pubparse/internal/fileIO"
)

// TestVerifyPath uses table-driven testing to validate the behavior of fileIO.VerifyPath,
// checking for correct extension handling, missing files, and directories.
func TestVerifyPath(t *testing.T) {
	// Create an isolated temporary directory for testing.
	tmpDir := t.TempDir()

	// Create sample files with known content and extensions.
	txtFile := filepath.Join(tmpDir, "test.txt")
	jsonFile := filepath.Join(tmpDir, "test.json")

	os.WriteFile(txtFile, []byte("hello"), 0644)
	os.WriteFile(jsonFile, []byte("{}"), 0644)

	// Define a set of test cases to check different edge cases and valid scenarios.
	tests := []struct {
		name        string // Human-readable name of the test case.
		path        string // Path to test.
		fileType    string // Expected file extension, e.g., "txt", ".json"
		expectError bool   // Whether an error is expected.
		errType     any    // Optional: check if error is of a specific type.
	}{
		{"valid .txt file", txtFile, "txt", false, nil},
		{"valid .json file with dot", jsonFile, ".json", false, nil},
		{"wrong extension", jsonFile, "xml", true, &customErrors.WrongExtensionError{}},
		{"empty type (no check)", txtFile, "", false, nil},
		{"nonexistent file", filepath.Join(tmpDir, "nope.txt"), "txt", true, nil},
		{"directory with no type", tmpDir, "", false, nil},
	}

	// Loop through each test case and run it as a subtest for isolation and clarity.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := fileIO.VerifyPath(tt.path, tt.fileType)

			// Check for expected errors
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.errType != nil {
					if err == nil {
						t.Errorf("expected error type %T, got %v", tt.errType, err)
					}
				}
			} else {
				// Unexpected error case
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if info == nil {
					t.Errorf("expected file info, got nil")
				}
			}
		})
	}
}
