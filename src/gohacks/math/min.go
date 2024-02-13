package math

import "golang.org/x/exp/constraints"

// Generic min function that should be in the Go standard library.
func Min[T constraints.Ordered](element1, element2 T) T {
	if element1 < element2 {
		return element1
	}
	return element2
}
