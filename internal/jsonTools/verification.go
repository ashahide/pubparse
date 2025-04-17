package jsonTools

import (
	"fmt"
	"log"

	"github.com/xeipuuv/gojsonschema"
)

//
// ------------------------ ValidateJsonAgainstSchema ------------------------
//

// ValidateJsonAgainstSchema checks whether a JSON file conforms to a JSON Schema definition.
//
// Arguments:
//   - path_to_json:   Path to the JSON file to be validated.
//   - path_to_schema: Path to the JSON Schema file (.json) to validate against.
//
// Returns:
//   - error: Any error encountered during validation (e.g., file access, schema parse).
//   - If the schema is invalid or the JSON does not match, the error is logged and returned.
//
// Behavior:
//   - Uses the github.com/xeipuuv/gojsonschema package for validation.
//   - Prints all validation errors if the JSON fails schema validation.
//   - Calls log.Fatal for fatal load/parse errors (can be replaced with return for safer handling).
func ValidateJsonAgainstSchema(path_to_json string, path_to_schema string) error {
	// Load the schema and document from file paths using the file:// URI scheme
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + path_to_schema)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + path_to_json)

	// Perform schema validation
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		// Fatal error: likely unable to load or parse one of the files
		log.Fatal("Validation error:", err)
	}

	// Output validation results
	if result.Valid() {
		fmt.Println("JSON is valid against the schema.")
	} else {
		fmt.Println("JSON is NOT valid against the schema.")
		for _, desc := range result.Errors() {
			fmt.Println(" -", desc)
		}
	}

	// Return error (nil if valid, or validation failure)
	return err
}
