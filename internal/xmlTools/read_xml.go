package xmlTools

import (
	"encoding/xml"
	"fmt"
	"os"
)

// ------------------------ ParsePubmedXML ------------------------

/*
ParsePubmedXML attempts to detect and parse a PubMed or PMC XML file.

It tries parsing the file into one of the following known formats:
  - *PubmedArticleSet
  - *PubmedBookArticleSet
  - *PMCArticle

Parameters:
  - filePath: Path to the XML file on disk.

Returns:
  - Parsed result as an interface{} (e.g., *PubmedArticleSet, *PubmedBookArticleSet, or *PMCArticle).
  - Error if the file cannot be read or parsed as a recognized structure.

Behavior:
  - Reads the XML file into memory.
  - Attempts to unmarshal into PubmedArticleSet.
  - If no articles are found, attempts PubmedBookArticleSet.
  - Then tries PMCArticle.
  - Returns an error if none of the known formats match.
*/
func ParsePubmedXML(filePath string) (interface{}, error) {
	xmlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Attempt to parse as PubmedArticleSet
	var articleSet PubmedArticleSet
	if err := xml.Unmarshal(xmlBytes, &articleSet); err == nil && len(articleSet.PubmedArticles) > 0 {
		return &articleSet, nil
	}

	// Attempt to parse as PubmedBookArticleSet
	var bookSet PubmedBookArticleSet
	if err := xml.Unmarshal(xmlBytes, &bookSet); err == nil && len(bookSet.PubmedBookArticles) > 0 {
		return &bookSet, nil
	}

	// Attempt to parse as PMCArticle
	var pmc PMCArticle
	if err := xml.Unmarshal(xmlBytes, &pmc); err != nil {
		fmt.Println("PMC unmarshal error:", err)
	} else if pmc.XMLName.Local != "article" {
		fmt.Println("PMC parsed, but root element was:", pmc.XMLName.Local)
	} else {
		return &pmc, nil
	}

	return nil, fmt.Errorf("unrecognized PubMed or PMC XML structure")
}

// ------------------------ NormalizePubmedArticleSet ------------------------

/*
NormalizePubmedArticleSet ensures that slice fields inside a PubmedArticleSet
are initialized to empty slices instead of nil.

This prevents fields from being serialized as `null` in JSON output,
which is important for JSON Schema validation.

Parameters:
  - data: Parsed PubmedArticleSet (interface{}).

Behavior:
  - If data is a *PubmedArticleSet:
  - Ensures MedlineCitation.KeywordList is an empty []string if nil.
  - Ensures PubmedData.ReferenceList is an empty []Reference if nil.
  - Ensures Unknown is an empty []UnknownElement if nil.

Note:
  - PubmedBookArticleSet and PMCArticle normalization can be added separately.
*/
func NormalizePubmedArticleSet(data interface{}) {
	switch v := data.(type) {
	case *PubmedArticleSet:
		for i := range v.PubmedArticles {
			article := &v.PubmedArticles[i]

			// Ensure KeywordList is non-nil
			if article.MedlineCitation.KeywordList == nil {
				article.MedlineCitation.KeywordList = []string{}
			}

			// Ensure ReferenceList is non-nil
			if article.PubmedData.ReferenceList == nil {
				article.PubmedData.ReferenceList = []Reference{}
			}

			// Ensure Unknown is non-nil
			if article.Unknown == nil {
				article.Unknown = []UnknownElement{}
			}
		}
	}
}

// ------------------------ NormalizePMCArticle ------------------------

/*
NormalizePMCArticle ensures that fields inside a PMCArticle
are initialized to valid empty types instead of being nil.

This guarantees consistency for JSON serialization and schema validation.

Parameters:
  - article: Pointer to a PMCArticle to normalize.

Behavior:
  - Initializes missing FloatsGroup, Back, Body sections.
  - Ensures paragraphs and references are initialized to empty slices.
  - Initializes missing affiliation information inside contributors.
*/
func NormalizePMCArticle(article *PMCArticle) {
	// Ensure FloatsGroup is non-nil
	if article.FloatsGroup == nil {
		article.FloatsGroup = &PMCFloatsGroup{Figures: []PMCFigure{}}
	}

	// Ensure Back is non-nil and subfields are initialized
	if article.Back == nil {
		article.Back = &PMCBack{
			Acknowledgments: &PMCAcknowledgments{Paragraphs: []string{}},
			References:      &PMCReferences{References: []PMCReference{}},
		}
	} else {
		if article.Back.Acknowledgments == nil {
			article.Back.Acknowledgments = &PMCAcknowledgments{Paragraphs: []string{}}
		}
		if article.Back.References == nil {
			article.Back.References = &PMCReferences{References: []PMCReference{}}
		}
	}

	// Ensure Body is non-nil and all sections have initialized paragraphs
	if article.Body == nil {
		article.Body = &PMCBody{Sections: []PMCSection{}}
	} else {
		for i := range article.Body.Sections {
			if article.Body.Sections[i].Paragraphs == nil {
				article.Body.Sections[i].Paragraphs = []string{}
			}
		}
	}

	// Ensure all contributors have a non-nil affiliation (Aff)
	for g := range article.Front.ArticleMeta.ContribGroup {
		for c := range article.Front.ArticleMeta.ContribGroup[g].Contrib {
			if article.Front.ArticleMeta.ContribGroup[g].Contrib[c].Aff == nil {
				article.Front.ArticleMeta.ContribGroup[g].Contrib[c].Aff = &PMCAff{}
			}
		}
	}

	// Ensure Abstract has initialized paragraphs
	if article.Front.ArticleMeta.Abstract != nil && article.Front.ArticleMeta.Abstract.Paragraphs == nil {
		article.Front.ArticleMeta.Abstract.Paragraphs = []string{}
	}
}
