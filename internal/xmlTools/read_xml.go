package xmlTools

import (
	"encoding/xml"
	"fmt"
	"os"
)

//
// ------------------------ ParsePubmedXML ------------------------
//

// ParsePubmedXML attempts to detect and parse a PubMed XML file as either a
// PubmedArticleSet or PubmedBookArticleSet.
//
// Parameters:
//   - filePath: Path to the XML file on disk.
//
// Returns:
//   - Parsed result as an interface{} (either *PubmedArticleSet or *PubmedBookArticleSet).
//   - Error if the file cannot be read or parsed as a known PubMed format.
//
// Behavior:
//   - Reads the XML file into memory.
//   - Tries decoding into PubmedArticleSet first.
//   - If no articles found, tries decoding as PubmedBookArticleSet.
//   - If neither succeeds, returns an error indicating unknown structure.
func ParsePubmedXML(filePath string) (interface{}, error) {
	xmlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var articleSet PubmedArticleSet
	if err := xml.Unmarshal(xmlBytes, &articleSet); err == nil && len(articleSet.PubmedArticles) > 0 {
		return &articleSet, nil
	}

	var bookSet PubmedBookArticleSet
	if err := xml.Unmarshal(xmlBytes, &bookSet); err == nil && len(bookSet.PubmedBookArticles) > 0 {
		return &bookSet, nil
	}

	return nil, fmt.Errorf("unrecognized PubMed XML structure")
}

//
// ------------------------ NormalizePubmedArticleSet ------------------------
//

// NormalizePubmedArticleSet ensures that key slice fields inside a PubmedArticleSet
// are never nil, which prevents them from being serialized as `null` in JSON output.
//
// This is especially important for JSON Schema validation, which expects arrays like
// `KeywordList` and `ReferenceList` to be present (even if empty).
//
// Parameters:
//   - data: The input PubmedArticleSet or PubmedBookArticleSet as an interface{}.
//
// Behavior:
//   - If data is a *PubmedArticleSet:
//   - Iterates over all PubmedArticles.
//   - Ensures MedlineCitation.KeywordList is an empty []string if nil.
//   - Ensures PubmedData.ReferenceList is an empty []Reference if nil.
//   - Ensures Unknown is an empty []UnknownElement if nil.
//   - Book normalization logic can be added in the future.
func NormalizePubmedArticleSet(data interface{}) {
	switch v := data.(type) {
	case *PubmedArticleSet:
		for i := range v.PubmedArticles {
			article := &v.PubmedArticles[i]

			// Normalize KeywordList
			if article.MedlineCitation.KeywordList == nil {
				article.MedlineCitation.KeywordList = []string{}
			}

			// Normalize ReferenceList
			if article.PubmedData.ReferenceList == nil {
				article.PubmedData.ReferenceList = []Reference{}
			}

			// Normalize Unknown
			if article.Unknown == nil {
				article.Unknown = []UnknownElement{}
			}
		}
	}
}
