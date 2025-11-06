package gohacks

import (
	"fmt"
)

// WrapIfError wraps the given error in the given text.
func WrapIfError(err error, text string, a ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(text, append(a, err)...)
}
