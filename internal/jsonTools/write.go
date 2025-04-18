package jsonTools

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ashahide/pubparse/internal/fileIO"
	"github.com/ashahide/pubparse/internal/makeReports"
	"github.com/ashahide/pubparse/internal/xmlTools"
)

//
// ------------------------ ConvertToJSON ------------------------
//

/*
ConvertToJSON serializes a normalized PubMed or PMC structure to JSON and validates it against a schema.

Parameters:
  - result: The parsed and normalized data structure. Must be one of:
    *xmlTools.PubmedArticleSet, *xmlTools.PubmedBookArticleSet, or *xmlTools.PMCArticle.
  - fileName: Path to save the output JSON.
  - schemaPath: Path to the JSON Schema file to validate against.

Behavior:
  - Serializes the structure into compact JSON using encoding/json.
  - Writes the output to the specified file.
  - Validates the output file against the given schema using ValidateJsonAgainstSchema.

Fatal:
  - Terminates the program via log.Fatal if marshaling or file writing fails.

Returns:
  - An error from schema validation if validation fails; otherwise nil.
*/
func ConvertToJSON(result interface{}, fileName, schemaPath string) error {
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal("failed to marshal result to JSON:", err)
	}

	if err := os.WriteFile(fileName, jsonData, 0644); err != nil {
		log.Fatal("failed to write JSON to file:", err)
	}

	return ValidateJsonAgainstSchema(fileName, schemaPath)
}

//
// ------------------------ serializeAndValidate ------------------------
//

/*
serializeAndValidate prepares and saves a parsed PubMed or PMC structure to JSON and validates it.

Parameters:
  - data: The parsed structure, such as *PubmedArticleSet, *PubmedBookArticleSet, or *PMCArticle.
  - outputPath: The full path where the JSON should be written.

Behavior:
  - Normalizes the data depending on its type.
  - Selects the appropriate JSON schema.
  - Calls ConvertToJSON to write and validate the file.

Returns:
  - The path to the schema used.
  - An error if serialization or validation fails.
*/
func serializeAndValidate(data interface{}, outputPath string) (string, error) {
	switch v := data.(type) {
	case *xmlTools.PubmedArticleSet:
		xmlTools.NormalizePubmedArticleSet(v)
		schema := filepath.Join("internal", "jsonTools", "pubmed_json_schema.json")
		return schema, ConvertToJSON(v, outputPath, schema)

	case *xmlTools.PubmedBookArticleSet:
		xmlTools.NormalizePubmedArticleSet(v)
		schema := filepath.Join("internal", "jsonTools", "pubmed_json_schema.json")
		return schema, ConvertToJSON(v, outputPath, schema)

	case *xmlTools.PMCArticle:
		xmlTools.NormalizePMCArticle(v)
		schema := filepath.Join("internal", "jsonTools", "pmc_json_schema.json")
		return schema, ConvertToJSON(v, outputPath, schema)

	default:
		return "", fmt.Errorf("unsupported data type for serialization")
	}
}

//
// ------------------------ processFile ------------------------
//

/*
processFile executes the full processing pipeline for a single input file:
parse XML, normalize, convert to JSON, validate, and log result.

Parameters:
  - i: Index of the file in the file list.
  - args: The input/output file path configuration.
  - mode: "pubmed" or "pmc", used to guide XML parsing.
  - report: Open report file handle for logging.
  - mu: Mutex to ensure thread-safe access to the report file.
  - start: Start time of the entire processing batch (for progress).
  - doneCounter: Atomic counter tracking how many files have been processed.

Returns:
  - An error if any stage in processing fails; otherwise nil.
*/
func processFile(
	i int,
	args fileIO.Arguments,
	mode string,
	report *os.File,
	mu *sync.Mutex,
	start time.Time,
	doneCounter *int32,
) error {
	fin := args.InputPath.Files[i]
	fout := args.OutputPath.Files[i]

	// Ensure output file is created before writing
	if err := fileIO.MakeFile(fout); err != nil {
		return fmt.Errorf("failed to create output file %q: %w", fout, err)
	}

	// Parse XML file into appropriate structure
	data, err := xmlTools.ParsePubmedXML(fin)
	if err != nil {
		return fmt.Errorf("failed to parse XML %q: %w", fin, err)
	}

	// Convert and validate
	if _, convErr := serializeAndValidate(data, fout); convErr != nil {
		return fmt.Errorf("failed to convert to JSON for %q: %w", fout, convErr)
	}

	// Write mapping to report
	if report != nil {
		if err := makeReports.WriteToReport(report, mu, fin, fout); err != nil {
			return fmt.Errorf("failed to write to report: %w", err)
		}
	}

	// Increment progress
	atomic.AddInt32(doneCounter, 1)
	return nil
}

//
// ------------------------ ProcessAllFiles ------------------------
//

/*
ProcessAllFiles runs the main pipeline on all input files with parallel workers.

Parameters:
  - args: Holds the input/output file lists and paths.
  - mode: Either "pubmed" or "pmc", used to guide XML parsing logic.
  - report: Open file handle to write report log entries.
  - workers: Number of parallel goroutines to spawn for concurrent file processing.

Behavior:
  - Spawns a progress tracker in the background.
  - Uses a semaphore (channel) to enforce the worker limit.
  - Each file is processed in its own goroutine, with error handling and safe report logging.
  - All workers wait using a sync.WaitGroup.

Returns:
  - The first encountered error during processing, or nil if all files succeed.
*/
func ProcessAllFiles(args fileIO.Arguments, mode string, report *os.File, workers int) error {
	startTime := time.Now()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var doneCount int32
	stopCh := make(chan struct{})

	// Start progress tracker in background
	go makeReports.TrackProgress(len(args.InputPath.Files), &doneCount, startTime, stopCh)

	sema := make(chan struct{}, workers)                   // concurrency limiter
	errChan := make(chan error, len(args.InputPath.Files)) // error collection channel

	// Launch worker goroutines
	for i := range args.InputPath.Files {
		wg.Add(1)
		sema <- struct{}{} // acquire slot

		go func(i int) {
			defer wg.Done()
			defer func() { <-sema }() // release slot

			if err := processFile(i, args, mode, report, &mu, startTime, &doneCount); err != nil {
				errChan <- err
			}
		}(i)
	}

	// Wait for all processing to complete
	wg.Wait()
	close(errChan)
	close(stopCh)

	// Return first encountered error, if any
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
