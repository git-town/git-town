package util

// Ensure asserts that the given condition is true.
// If not, it ends the application with the given error message.
func Ensure(condition bool, error string) {
	if !condition {
		ExitWithErrorMessage(error)
	}
}
