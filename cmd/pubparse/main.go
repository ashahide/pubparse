package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ashahide/pubparse/internal/fileIO"
)

func main() {
	if err := run(); err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
}

func run() error {
	var args fileIO.Arguments

	// Define command line flags
	flag.StringVar(&args.InputPath.Path, "i", "", "Path to the input file or directory of files")
	flag.Parse()

	// Process Inputs
	fileIO.HandleInputs(&args)

	// Print file info
	fmt.Println("Input Path:", args.InputPath.Path)
	for _, f := range args.InputPath.Files {
		fmt.Printf(" - %s (%d bytes)\n", f.Name(), f.Size())
	}

	fmt.Println("Output Path:", args.OutputPath.Path)
	for _, f := range args.OutputPath.Files {
		fmt.Printf(" - %s (%d bytes)\n", f.Name(), f.Size())
	}

	return nil
}
