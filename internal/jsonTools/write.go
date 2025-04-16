package jsonTools

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/ashahide/pubparse/internal/xmlTools"
)

func ConvertToJson(result *xmlTools.PubmedArticleSet, output_path string, file_name os.FileInfo) error {

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(output_path, file_name.Name()), jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Validate the JSON output against the schema
	schemaPath := filepath.Join("internal", "jsonTools", "pubmed_json_schema.json")
	err = ValidateJsonAgainstSchema(filepath.Join(output_path, file_name.Name()), schemaPath)

	return err

}
