package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ashahide/pubparse/internal/fileIO"
)

func ParseArgs() (*fileIO.Arguments, string, error) {
	if len(os.Args) < 2 {
		return nil, "", errors.New("usage: pubparse [pubmed|pmc] -i input -o output")
	}

	var args fileIO.Arguments
	var mode string

	switch os.Args[1] {
	case "pubmed":
		mode = "pubmed"
		pubmedCmd := flag.NewFlagSet("pubmed", flag.ExitOnError)
		input := pubmedCmd.String("i", "", "Input file or directory")
		output := pubmedCmd.String("o", "", "Output file or directory")
		_ = pubmedCmd.Parse(os.Args[2:])

		args.InputPath.Path = *input
		args.OutputPath.Path = *output

	case "pmc":
		mode = "pmc"
		pmcCmd := flag.NewFlagSet("pmc", flag.ExitOnError)
		input := pmcCmd.String("i", "", "PMC input path")
		output := pmcCmd.String("o", "", "PMC output path")
		_ = pmcCmd.Parse(os.Args[2:])

		args.InputPath.Path = *input
		args.OutputPath.Path = *output

	default:
		return nil, "", fmt.Errorf("unknown subcommand: %s", os.Args[1])
	}

	return &args, mode, nil
}
