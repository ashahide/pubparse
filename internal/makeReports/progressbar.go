package makeReports

import (
	"fmt"
	"strings"
	"time"
)

//
// ------------------------ printProgressBar ------------------------
//

/*
printProgressBar prints a simple progress bar to the terminal based on current progress.

Parameters:
  - current: The current iteration or progress count.
  - total: The total number of iterations or target value.

Behavior:
  - Calculates the proportion of work completed.
  - Renders a fixed-width progress bar (50 characters wide).
  - Displays the percentage complete as a number to the right.
  - Uses carriage return (\r) to overwrite the current line in the terminal.
*/
func PrintProgressBar(current, total int, startTime time.Time) {
	width := 50
	progress := int(float64(current) / float64(total) * float64(width))
	bar := "[" + Repeat("=", progress) + Repeat(" ", width-progress) + "]"
	percent := float64(current) / float64(total) * 100
	elapsed := time.Since(startTime)
	elapsedTime := fmt.Sprintf("%02d:%02d:%02d", int(elapsed.Hours()), int(elapsed.Minutes())%60, int(elapsed.Seconds())%60)
	fmt.Printf("\r%s %.0f%% (%d/%d) - Elapsed Time: %s", bar, percent, current, total, elapsedTime)
}

//
// ------------------------ repeat ------------------------
//

/*
repeat returns a string consisting of `count` copies of the first rune in `char`.

Parameters:
  - char: A string whose first character will be repeated.
  - count: Number of times to repeat the character.

Returns:
  - A string made up of `count` copies of the first rune in `char`.

Note:
  - Only the first character of `char` is used. If `char` is empty, this returns an empty string.
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
