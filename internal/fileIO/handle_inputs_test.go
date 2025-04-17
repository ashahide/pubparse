package fileIO_test

// import (
// 	"os"
// 	"path/filepath"
// 	"testing"

// 	"github.com/ashahide/pubparse/internal/fileIO"
// )

// // Test single XML file input
// func TestHandleInputs_SingleFile(t *testing.T) {
// 	tempDir := t.TempDir()
// 	inputFile := filepath.Join(tempDir, "test_input.xml")
// 	os.WriteFile(inputFile, []byte(`<root></root>`), 0644)

// 	args := &fileIO.Arguments{
// 		InputPath: fileIO.PathInfo{
// 			Path: inputFile,
// 		},
// 	}

// 	err := fileIO.HandleInputs(args)
// 	if err != nil {
// 		t.Fatalf("HandleInputs failed: %v", err)
// 	}

// 	if args.InputPath.Info == nil {
// 		t.Fatal("InputPath.Info is nil")
// 	}
// 	if len(args.InputPath.Files) != 1 {
// 		t.Fatalf("expected 1 input file, got %d", len(args.InputPath.Files))
// 	}
// 	if len(args.OutputPath.Files) != 1 {
// 		t.Fatalf("expected 1 output file, got %d", len(args.OutputPath.Files))
// 	}
// 	if filepath.Ext(args.OutputPath.Files[0].Name()) != ".json" {
// 		t.Errorf("expected .json extension, got %s", args.OutputPath.Files[0].Name())
// 	}
// }

// // Test multiple XML files in a directory
// func TestHandleInputs_Directory(t *testing.T) {
// 	tempDir := t.TempDir()
// 	fileNames := []string{"a.xml", "b.xml", "c.xml"}
// 	for _, name := range fileNames {
// 		os.WriteFile(filepath.Join(tempDir, name), []byte(`<doc></doc>`), 0644)
// 	}

// 	args := &fileIO.Arguments{
// 		InputPath: fileIO.PathInfo{
// 			Path: tempDir,
// 		},
// 	}

// 	err := fileIO.HandleInputs(args)
// 	if err != nil {
// 		t.Fatalf("HandleInputs failed: %v", err)
// 	}

// 	if len(args.InputPath.Files) != len(fileNames) {
// 		t.Errorf("expected %d input files, got %d", len(fileNames), len(args.InputPath.Files))
// 	}
// 	if len(args.OutputPath.Files) != len(fileNames) {
// 		t.Errorf("expected %d output files, got %d", len(fileNames), len(args.OutputPath.Files))
// 	}
// 	for _, out := range args.OutputPath.Files {
// 		if filepath.Ext(out.Name()) != ".json" {
// 			t.Errorf("expected .json extension, got %s", out.Name())
// 		}
// 	}
// }

// // Test non-existent path
// func TestHandleInputs_InvalidPath(t *testing.T) {
// 	args := &fileIO.Arguments{
// 		InputPath: fileIO.PathInfo{
// 			Path: "/nonexistent/path/to/file.xml",
// 		},
// 	}

// 	err := fileIO.HandleInputs(args)
// 	if err == nil {
// 		t.Fatal("expected error for non-existent path, got nil")
// 	}
// }

// // Test file with wrong extension is skipped
// func TestHandleInputs_SkipWrongExtension(t *testing.T) {
// 	tempDir := t.TempDir()
// 	_ = os.WriteFile(filepath.Join(tempDir, "skip.txt"), []byte(`text`), 0644)
// 	_ = os.WriteFile(filepath.Join(tempDir, "valid.xml"), []byte(`<root/>`), 0644)

// 	args := &fileIO.Arguments{
// 		InputPath: fileIO.PathInfo{
// 			Path: tempDir,
// 		},
// 	}

// 	err := fileIO.HandleInputs(args)
// 	if err != nil {
// 		t.Fatalf("HandleInputs failed: %v", err)
// 	}

// 	if len(args.InputPath.Files) != 1 {
// 		t.Errorf("expected 1 valid input file, got %d", len(args.InputPath.Files))
// 	}
// 	if args.InputPath.Files[0].Name() != "valid.xml" {
// 		t.Errorf("expected input file 'valid.xml', got %s", args.InputPath.Files[0].Name())
// 	}
// }
