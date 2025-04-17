package jsonTools

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/ashahide/pubparse/internal/xmlTools"
)

//
// ------------------------ ConvertToJson ------------------------
//

// ConvertToJson serializes a parsed PubMed XML data structure into a formatted JSON file
// and validates the resulting output against a predefined JSON Schema.
//
// Arguments:
//   - result: A pointer to a PubmedArticleSet (the parsed XML result).
//   - file_name: The full file path where the output .json file should be written.
//
// Behavior:
//   - Uses json.MarshalIndent for pretty-printed output.
//   - Writes the JSON file to the specified location with mode 0644.
//   - Validates the file against pubmed_json_schema.json.
//
// Returns:
//   - error: If marshaling, writing, or schema validation fails.
//     (If marshaling or writing fails, log.Fatal is called instead.)
func ConvertToJson(result *xmlTools.PubmedArticleSet, file_name string) error {
	// Marshal the result into a human-readable, indented JSON format
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		// Fatal: stop the program on JSON marshal error
		log.Fatal(err)
	}

	// Write the JSON data to the specified output file
	err = os.WriteFile(file_name, jsonData, 0644)
	if err != nil {
		// Fatal: stop the program on write error
		log.Fatal(err)
	}

	// Construct the absolute path to the validation schema
	schemaPath := filepath.Join("internal", "jsonTools", "pubmed_json_schema.json")

	// Run schema validation on the newly created JSON file
	err = ValidateJsonAgainstSchema(file_name, schemaPath)

	// Return any validation error (or nil if validation passes)
	return err
}
