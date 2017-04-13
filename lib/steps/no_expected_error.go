package steps

// NoExpectedError is a partial Step implementation used for composition.
type NoExpectedError struct {
	NoAbortStep
	NoAutomaticAbortOnError
	NoContinueStep
}
