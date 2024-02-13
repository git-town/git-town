package math

// Generic max function that should be in the Go standard library.
func Max(element1, element2 int) int {
	if element1 > element2 {
		return element1
	}
	return element2
}
