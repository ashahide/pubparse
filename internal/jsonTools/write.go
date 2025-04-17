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

// ConvertToJson serializes a parsed PubMed XML structure—either a PubmedArticleSet
// or PubmedBookArticleSet—into formatted JSON and validates the output against a JSON Schema.
//
// Parameters:
//   - result: The parsed XML data as an interface{}, expected to be either
//     *xmlTools.PubmedArticleSet or *xmlTools.PubmedBookArticleSet.
//   - file_name: Full destination path where the output JSON should be saved.
//
// Behavior:
//   - Calls NormalizePubmedArticleSet to ensure critical array fields (e.g., KeywordList, ReferenceList)
//     are non-nil in PubmedArticleSet (to avoid `null` in JSON).
//   - Serializes the input structure to indented, human-readable JSON using json.MarshalIndent.
//   - Writes the serialized output to the specified file path.
//   - Validates the written JSON file against the pubmed_json_schema.json schema.
//
// Returns:
//   - error: nil if successful, or an error from schema validation.
//     Critical failures during marshalling or writing will abort execution via log.Fatal.
func ConvertToJson(result interface{}, file_name string) error {
	// Normalize known structures (currently only PubmedArticleSet).
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
