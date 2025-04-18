package jsonTools

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ashahide/pubparse/internal/xmlTools"
)

//
// ------------------------ ConvertToJson ------------------------
//

/*
ConvertToJson serializes a parsed PubMed XML structure (e.g., PubmedArticleSet or PubmedBookArticleSet)
into a pretty-printed JSON file and validates it against a JSON Schema.

Parameters:
  - result: Parsed XML structure, expected to be one of:
    *xmlTools.PubmedArticleSet or *xmlTools.PubmedBookArticleSet.
  - file_name: Full path where the JSON output should be saved.
  - schemapath: Path to the JSON schema used for validation.

Behavior:
  - Normalizes known XML types (currently only PubmedArticleSet) to ensure required
    slice fields are not nil (e.g., KeywordList, ReferenceList).
  - Serializes the normalized structure into human-readable indented JSON.
  - Saves the serialized JSON to the specified file path.
  - Validates the written file against the JSON schema.

Returns:
  - error: Validation error if the JSON does not conform to the schema; otherwise nil.

Note:
  - Any failure during JSON encoding or file writing is fatal and will terminate the program.
*/
func ConvertToJson(result interface{}, file_name string, schemapath string) error {
	// Normalize nil slices in known types (avoids `null` in JSON).
	xmlTools.NormalizePubmedArticleSet(result)

	// Marshal the data to indented JSON.
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal result to JSON:", err)
	}

	// Write the JSON to the specified file.
	err = os.WriteFile(file_name, jsonData, 0644)
	if err != nil {
		log.Fatal("failed to write JSON to file:", err)
	}

	// Validate the written file against the provided JSON schema.
	err = ValidateJsonAgainstSchema(file_name, schemapath)

	// Return validation result (or nil if successful).
	return err
}
