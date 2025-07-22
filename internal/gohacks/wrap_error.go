package gohacks

import "fmt"

func WrapError(err error, text string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(text, err)
}
