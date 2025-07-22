package gohacks

import "fmt"

// If the given error is nil, returns nil.
// If the given error is not nil, returns it wrapped in the given error message.
func WrapIfError(err error, text string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(text, err)
}
