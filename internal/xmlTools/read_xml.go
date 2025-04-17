package xmlTools

import (
	"encoding/xml"
	"fmt"
	"os"
)

//
// ------------------------ ParsePubmedArticleSet ------------------------
//

// ParsePubmedArticleSet opens and parses a PubMed XML file into a PubmedArticleSet struct.
//
// Arguments:
//   - filePath: Full path to the XML file to be parsed.
//
// Returns:
//   - *PubmedArticleSet: A pointer to the parsed struct containing articles and metadata.
//   - error: If the file cannot be opened or if XML parsing fails.
//
// Behavior:
//   - Opens the file using os.Open and ensures it is closed with defer.
//   - Uses encoding/xml's Decoder to parse the XML stream.
//   - Returns a descriptive error if parsing fails.
func ParsePubmedArticleSet(filePath string) (*PubmedArticleSet, error) {
	// Attempt to open the specified XML file
	f, err := os.Open(filePath)
	if err != nil {
		// Return a wrapped error if the file cannot be opened
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// Declare a variable to store the parsed XML structure
	var articleSet PubmedArticleSet

	// Create a streaming XML decoder for efficient parsing
	decoder := xml.NewDecoder(f)

	// Attempt to decode the entire XML file into the articleSet structure
	if err := decoder.Decode(&articleSet); err != nil {
		// Return a wrapped error if parsing fails
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	// Return the successfully parsed result
	return &articleSet, nil
}
