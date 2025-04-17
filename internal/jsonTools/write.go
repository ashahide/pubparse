package jsonTools

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/ashahide/pubparse/internal/xmlTools"
)

func ConvertToJson(result *xmlTools.PubmedArticleSet, file_name string) error {

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(file_name, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Validate the JSON output against the schema
	schemaPath := filepath.Join("internal", "jsonTools", "pubmed_json_schema.json")
	err = ValidateJsonAgainstSchema(file_name, schemaPath)

	return err

}
