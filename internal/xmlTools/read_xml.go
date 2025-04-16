package xmlTools

import (
	"encoding/xml"
	"fmt"
	"os"
)

func ParsePubmedArticleSet(filePath string) (*PubmedArticleSet, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var articleSet PubmedArticleSet
	decoder := xml.NewDecoder(f)
	if err := decoder.Decode(&articleSet); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	return &articleSet, nil
}
