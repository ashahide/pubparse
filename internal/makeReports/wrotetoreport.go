package makeReports

import (
	"fmt"
	"os"
	"sync"
)

//
// ------------------------ WriteToReport ------------------------
//

/*
WriteToReport writes a log entry to the report file, mapping an input file to its output,
in a thread-safe manner using a mutex.

Parameters:
  - report: An open *os.File for writing report entries.
  - mu: Pointer to a sync.Mutex used to guard concurrent access to the file.
  - fin: Path to the input XML file.
  - fout: Path to the output JSON file.

Behavior:
  - Acquires a lock on the mutex before performing any file operations.
  - Appends a formatted line indicating which input file was converted to which output.
  - Forces a disk flush with Sync() to ensure the report is immediately saved.

Returns:
  - An error if writing or syncing the report file fails; otherwise nil.
*/
func WriteToReport(report *os.File, mu *sync.Mutex, fin, fout string) error {
	mu.Lock()         // Prevent concurrent writes from multiple goroutines
	defer mu.Unlock() // Always release lock, even on error

	// Format and append the input/output mapping
	if _, err := report.WriteString(fmt.Sprintf("\n>>> Input file: %s\t Output file: %s\n", fin, fout)); err != nil {
		return err
	}

	// Flush to disk to ensure durability
	return report.Sync()
}
