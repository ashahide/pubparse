package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ashahide/pubparse/internal/fileIO"
	"github.com/ashahide/pubparse/internal/jsonTools"
	"github.com/ashahide/pubparse/internal/xmlTools"
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

	for i := range args.InputPath.Files {
		fin := args.InputPath.Files[i]
		fout := args.OutputPath.Files[i]

		fmt.Println("Processing file:", fin.Name())

		result, err := xmlTools.ParsePubmedArticleSet(filepath.Join(args.InputPath.Path, fin.Name()))
		if err != nil {
			log.Fatal(err)
		}

		err = jsonTools.ConvertToJson(result, args.OutputPath.Path, fout)
		if err != nil {
			log.Fatal(err)
		}

	}

	return nil
}
