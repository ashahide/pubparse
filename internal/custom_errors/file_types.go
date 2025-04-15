package custom_errors

import (
	"fmt"
)

// WrongExtensionError represents an error due to file extension mismatch.
type WrongExtensionError struct {
	Expected string
	Actual   string
}

func (e *WrongExtensionError) Error() string {
	return fmt.Sprintf("wrong file extension: expected %s, got %s", e.Expected, e.Actual)
}
