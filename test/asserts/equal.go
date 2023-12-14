package asserts

import "fmt"

func Equal[T comparable](want, have T, msg string) error {
	if have == want {
		return nil
	}
	return fmt.Errorf(msg, have, want)
}
