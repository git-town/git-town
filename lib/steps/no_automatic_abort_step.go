package steps

type NoAutomaticAbortOnError struct{}

func (step NoAutomaticAbortOnError) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step NoAutomaticAbortOnError) ShouldAutomaticallyAbortOnError() bool {
	return false
}
