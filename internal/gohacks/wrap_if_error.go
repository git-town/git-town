package gohacks

import (
	"fmt"
)

// If the given error is nil, returns nil.
// If the given error is not nil, returns it wrapped in the given error message.
func WrapIfError(err error, text string, a ...any) error {
	if err == nil {
		return nil
	}
	args := []any{text}
	args = append(args, a...)
	args = append(args, err)
	return fmt.Errorf(text, args...)
}
