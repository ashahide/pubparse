package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ashahide/pubparse/internal/fileIO"
	"github.com/ashahide/pubparse/internal/jsonTools"
)

//
// ------------------------ main ------------------------
//

/*
main is the entry point for the pubparse command-line tool.

Behavior:
  - Delegates execution to run().
  - Logs any error returned and exits with status code 1.
*/
func main() {
	if err := run(); err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
}

//
// ------------------------ run ------------------------
//

/*
run handles all stages of execution: parsing CLI arguments, validating file paths,
logging metadata, and launching parallel XML-to-JSON conversion.

Returns:
  - An error if any stage in the processing pipeline fails.

Behavior:
  - Supports subcommands: "pubmed" or "pmc".
  - Required flags: -i (input), -o (output).
  - Optional flag: --workers (number of concurrent goroutines, default 8).
  - Validates file count alignment between input/output.
  - Creates a report file and logs session metadata.
  - Delegates parallel file processing to jsonTools.ProcessAllFiles.
*/
func run() error {
	// Ensure a valid subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: pubparse [pubmed|pmc] -i input_path -o output_path [--workers N]")
		os.Exit(1)
	}

	mode := os.Args[1] // Subcommand: "pubmed" or "pmc"
	var args fileIO.Arguments
	var workers int

	// Parse flags for the chosen mode
	switch mode {
	case "pubmed", "pmc":
		cmd := flag.NewFlagSet(mode, flag.ExitOnError)
		cmd.StringVar(&args.InputPath.Path, "i", "", "Path to the input file or directory")
		cmd.StringVar(&args.OutputPath.Path, "o", "", "Path to the output file or directory")
		cmd.IntVar(&workers, "workers", 8, "Number of concurrent workers (default: 8)")
		if err := cmd.Parse(os.Args[2:]); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown subcommand: %s\nUsage: pubparse [pubmed|pmc] -i input -o output [--workers N]", mode)
	}

	// Validate worker count
	if workers <= 0 {
		return fmt.Errorf("invalid number of workers: %d", workers)
	} else if workers > runtime.NumCPU() {
		fmt.Printf("Warning: Specified %d workers, but only %d CPU cores available. Setting workers = %d\n", workers, runtime.NumCPU(), runtime.NumCPU())
	}

	// Validate and resolve input/output paths and match file counts
	if err := fileIO.HandleInputs(&args); err != nil {
		return fmt.Errorf("input handling failed: %w", err)
	}
	if err := fileIO.HandleOutputs(&args); err != nil {
		return fmt.Errorf("output handling failed: %w", err)
	}
	if len(args.InputPath.Files) != len(args.OutputPath.Files) {
		return fmt.Errorf("input/output file count mismatch")
	}

	// Construct report file path and open it
	reportPath, err := filepath.Abs(filepath.Join(args.OutputPath.Path, "report.tsv"))
	if err != nil {
		return fmt.Errorf("failed to determine absolute report path: %w", err)
	}
	if err := fileIO.MakeFile(reportPath); err != nil {
		return fmt.Errorf("failed to create report file %q: %w", reportPath, err)
	}
	report, err := os.OpenFile(reportPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open report file %q: %w", reportPath, err)
	}
	defer report.Close()

	// Record metadata and start time in report
	startTime := time.Now()
	reportHeader := fmt.Sprintf(
		"\n>>> Starting Time: %s\n>>> Input Directory: %s\n>>> Output Directory: %s\n>>> Number of Inputs: %d\n>>> Workers: %d\n",
		startTime.Format("2006-01-02 15:04:05"),
		args.InputPath.Path,
		args.OutputPath.Path,
		len(args.InputPath.Files),
		workers,
	)
	if _, err := report.WriteString(reportHeader); err != nil {
		return fmt.Errorf("failed to write to report file %q: %w", reportPath, err)
	}

	// Echo metadata to console
	fmt.Println(">>> Input Path:", args.InputPath.Path)
	fmt.Println(">>> Output Path:", args.OutputPath.Path)
	fmt.Println(">>> Number of Inputs:", len(args.InputPath.Files))
	fmt.Println(">>> Workers:", workers)
	fmt.Println(">>> Starting Time:", startTime.Format("2006-01-02 15:04:05"))

	// Begin concurrent file processing
	if err := jsonTools.ProcessAllFiles(args, mode, report, workers); err != nil {
		return fmt.Errorf("processing failed: %w", err)
	}

	// Final summary
	fmt.Println("\n>>> Finished processing files.")
	fmt.Println(">>> Report file:", reportPath)
	fmt.Println(">>> Elapsed Time:", time.Since(startTime))
	fmt.Println(">>> Exiting...")
	return nil
}
