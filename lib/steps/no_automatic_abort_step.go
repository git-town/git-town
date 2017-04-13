package steps

// NoAutomaticAbortOnError is a partial Step implementation used for composition
type NoAutomaticAbortOnError struct{}

// GetAutomaticAbortErrorMessage returns the error message to display when this step
// cause the command to automatically abort.
func (step NoAutomaticAbortOnError) GetAutomaticAbortErrorMessage() string {
	return ""
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step NoAutomaticAbortOnError) ShouldAutomaticallyAbortOnError() bool {
	return false
}
