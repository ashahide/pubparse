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

// ConvertToJson serializes a parsed PubMed XML article set into formatted JSON
// and validates the result against a predefined JSON Schema.
//
// Parameters:
//   - result: A pointer to a PubmedArticleSet, representing the parsed PubMed XML structure.
//   - file_name: The full destination path where the output JSON file should be written.
//
// Behavior:
//   - Calls NormalizePubmedArticleSet to ensure all required JSON arrays are non-nil.
//   - Uses json.MarshalIndent to generate indented, human-readable output.
//   - Writes the JSON to the target file using os.WriteFile.
//   - Validates the written file against pubmed_json_schema.json.
//
// Returns:
//   - error: nil if successful, or an error from schema validation.
//     Note: Failures in marshalling or writing will cause the function to exit via log.Fatal.
func ConvertToJson(result *xmlTools.PubmedArticleSet, file_name string) error {
	// Ensure all slice fields are non-nil to avoid `null` values in the output JSON.
	xmlTools.NormalizePubmedArticleSet(result)

	// Serialize the data structure into formatted JSON.
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal result to JSON:", err)
	}

	// Write the resulting JSON to the output file.
	err = os.WriteFile(file_name, jsonData, 0644)
	if err != nil {
		log.Fatal("failed to write JSON to file:", err)
	}

	// Build the path to the validation schema.
	schemaPath := filepath.Join("internal", "jsonTools", "pubmed_json_schema.json")

	// Validate the output file against the schema.
	err = ValidateJsonAgainstSchema(file_name, schemaPath)

	// Return the result of validation (nil if success, error if invalid).
	return err
}
