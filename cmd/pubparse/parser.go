package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ashahide/pubparse/internal/fileIO"
)

//
// ------------------------ ParseArgs ------------------------
//

/*
ParseArgs parses and validates command-line arguments for the pubparse CLI tool.

Expected usage:

	pubparse [pubmed|pmc] -i input_path -o output_path [--workers N]

Supported subcommands:
  - pubmed: For parsing regular PubMed XML files.
  - pmc:    For parsing PMC full-text XML files.

Flags:

	-i: Input file or directory path.
	-o: Output file or directory path.
	--workers: (optional) Number of concurrent worker goroutines (default: 8).

Returns:
  - *fileIO.Arguments: Struct with resolved input/output paths.
  - string: Subcommand mode ("pubmed" or "pmc").
  - int: Number of workers requested.
  - error: Any error during argument parsing or validation.

Behavior:
  - Uses `flag.FlagSet` for mode-specific flag parsing.
  - Extracts input and output paths, and the number of workers.
  - Returns a formatted usage error if subcommand is missing or unrecognized.
*/
func ParseArgs() (*fileIO.Arguments, string, int, error) {
	if len(os.Args) < 2 {
		return nil, "", 0, errors.New("usage: pubparse [pubmed|pmc] -i input -o output [--workers N]")
	}

	var args fileIO.Arguments
	var mode string
	var workers int

	switch os.Args[1] {
	case "pubmed":
		mode = "pubmed"
		pubmedCmd := flag.NewFlagSet("pubmed", flag.ExitOnError)
		input := pubmedCmd.String("i", "", "Input file or directory")
		output := pubmedCmd.String("o", "", "Output file or directory")
		pubmedCmd.IntVar(&workers, "workers", 8, "Number of concurrent workers (default: 8)")

		if err := pubmedCmd.Parse(os.Args[2:]); err != nil {
			return nil, "", 0, err
		}
		args.InputPath.Path = *input
		args.OutputPath.Path = *output

	case "pmc":
		mode = "pmc"
		pmcCmd := flag.NewFlagSet("pmc", flag.ExitOnError)
		input := pmcCmd.String("i", "", "PMC input file or directory")
		output := pmcCmd.String("o", "", "PMC output file or directory")
		pmcCmd.IntVar(&workers, "workers", 8, "Number of concurrent workers (default: 8)")

		if err := pmcCmd.Parse(os.Args[2:]); err != nil {
			return nil, "", 0, err
		}
		args.InputPath.Path = *input
		args.OutputPath.Path = *output

	default:
		return nil, "", 0, fmt.Errorf("unknown subcommand: %s", os.Args[1])
	}

	return &args, mode, workers, nil
}
