package xmlTools

import (
	"encoding/xml"
	"fmt"
	"os"
)

//
// ------------------------ ParsePubmedArticleSet ------------------------
//

// ParsePubmedArticleSet reads a PubMed XML file and parses its content into a structured
// PubmedArticleSet object, which contains a slice of parsed PubMed articles.
//
// Parameters:
//   - filePath: The full path to the XML file to be opened and decoded.
//
// Returns:
//   - *PubmedArticleSet: A pointer to the resulting parsed data structure.
//   - error: An error if the file can't be opened or decoding fails.
//
// Behavior:
//   - Uses os.Open to read the file, and ensures the file is closed with defer.
//   - Uses encoding/xml.Decoder to stream and decode the XML content.
//   - Returns a fully populated PubmedArticleSet if parsing succeeds.
func ParsePubmedArticleSet(filePath string) (*PubmedArticleSet, error) {
	// Open the specified XML file
	f, err := os.Open(filePath)
	if err != nil {
		// Return an error if the file can't be opened
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close() // Ensure the file is closed after parsing

	// Prepare the target structure for XML decoding
	var articleSet PubmedArticleSet

	// Create an XML decoder for efficient token streaming
	decoder := xml.NewDecoder(f)

	// Decode the XML content into the target struct
	if err := decoder.Decode(&articleSet); err != nil {
		// Wrap and return decoding errors with context
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	// Return the fully parsed result
	return &articleSet, nil
}

//
// ------------------------ NormalizePubmedArticleSet ------------------------
//

// NormalizePubmedArticleSet ensures that key slice fields inside each PubmedArticle
// are never nil, which prevents them from being serialized as `null` in JSON output.
//
// This is especially important for JSON Schema validation, which expects arrays like
// `KeywordList` and `ReferenceList` to be present (even if empty).
//
// Parameters:
//   - set: A pointer to the PubmedArticleSet to normalize in-place.
//
// Behavior:
//   - Iterates over all PubmedArticles.
//   - Ensures MedlineCitation.KeywordList is initialized to an empty slice if nil.
//   - Ensures PubmedData.ReferenceList is initialized to an empty slice if nil.
func NormalizePubmedArticleSet(set *PubmedArticleSet) {
	for i := range set.PubmedArticles {
		a := &set.PubmedArticles[i]

		// Ensure KeywordList is an empty slice instead of nil
		if a.MedlineCitation.KeywordList == nil {
			a.MedlineCitation.KeywordList = []string{}
		}

		// Ensure ReferenceList is an empty slice instead of nil
		if a.PubmedData.ReferenceList == nil {
			a.PubmedData.ReferenceList = []Reference{}
		}
	}
}
