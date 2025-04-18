package makeReports

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

//
// ------------------------ PrintProgressBar ------------------------
//

/*
PrintProgressBar renders a terminal-based progress bar to show the percentage of tasks completed.

Parameters:
  - current: The current number of completed tasks.
  - total: The total number of tasks to complete.
  - startTime: The time at which processing began (used to display elapsed time).

Behavior:
  - Computes the progress percentage and elapsed time.
  - Displays a fixed-width progress bar (50 characters).
  - Uses carriage return (\r) to update the same line on each call.
*/
func PrintProgressBar(current, total int, startTime time.Time) {
	width := 50
	progress := int(float64(current) / float64(total) * float64(width))
	bar := "[" + Repeat("=", progress) + Repeat(" ", width-progress) + "]"
	percent := float64(current) / float64(total) * 100

	// Format elapsed time as HH:MM:SS
	elapsed := time.Since(startTime)
	elapsedTime := fmt.Sprintf("%02d:%02d:%02d", int(elapsed.Hours()), int(elapsed.Minutes())%60, int(elapsed.Seconds())%60)

	// Print dynamic progress line
	fmt.Printf("\r%s %.0f%% (%d/%d) - Elapsed Time: %s", bar, percent, current, total, elapsedTime)
}

//
// ------------------------ Repeat ------------------------
//

/*
Repeat returns a string made by repeating the first character of `char` exactly `count` times.

Parameters:
  - char: A string, from which only the first rune is repeated.
  - count: The number of repetitions.

Returns:
  - A string consisting of `count` copies of the first rune in `char`.

Notes:
  - Returns an empty string if `char` is empty or count is zero or negative.
  - Uses a strings.Builder for efficient memory allocation.
*/
func Repeat(char string, count int) string {
	if len(char) == 0 || count <= 0 {
		return ""
	}
	var b strings.Builder
	r := []rune(char)[0]
	for i := 0; i < count; i++ {
		b.WriteRune(r)
	}
	return b.String()
}

//
// ------------------------ TrackProgress ------------------------
//

/*
TrackProgress continuously updates the progress bar in a background goroutine.

Parameters:
  - total: The total number of tasks expected to complete.
  - done: Pointer to an atomic counter that tracks how many tasks have completed.
  - start: Timestamp indicating when the process began.
  - stop: Channel used to signal that processing is finished and the tracker should exit.

Behavior:
  - Starts a ticker that triggers every 100 milliseconds.
  - On each tick, it checks if the number of completed tasks (`done`) has changed.
  - If so, it updates the progress bar.
  - When the stop channel is closed, it prints the final 100% progress bar and exits.

Usage:
  - Should be run in a goroutine.
  - Complements concurrent processing pipelines by providing responsive visual feedback.
*/
func TrackProgress(total int, done *int32, start time.Time, stop <-chan struct{}) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	var lastPrinted int32 = -1

	for {
		select {
		case <-ticker.C:
			current := atomic.LoadInt32(done)
			if current != lastPrinted {
				PrintProgressBar(int(current), total, start)
				lastPrinted = current
			}

		case <-stop:
			// Print final state of the bar and move to a new line
			PrintProgressBar(total, total, start)
			fmt.Println()
			return
		}
	}
}
