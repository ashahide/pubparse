package jsonTools

import (
	"fmt"
	"log"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateJsonAgainstSchema(path_to_json string, path_to_schema string) error {

	// Load schema and data from files
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + path_to_json)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + path_to_schema)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		log.Fatal("Validation error:", err)
	}
	if result.Valid() {
		fmt.Println("JSON is valid against the schema.")
	} else {
		fmt.Println("JSON is NOT valid against the schema.")
		for _, desc := range result.Errors() {
			fmt.Println(" -", desc)
		}
	}

	return err
}
