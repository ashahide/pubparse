package fileIO

import (
	"testing"
)

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
			result := ChangeExtension(test.path, test.newExt)

			// Compare the actual result with the expected result
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}
