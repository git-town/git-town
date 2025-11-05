package gohacks

import (
	"fmt"
)

// WrapIfError returns nil if the given error is nil.
// If the given error is not nil, it returns it wrapped in the given error message.
func WrapIfError(err error, text string, a ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(text, append(a, err)...)
}
