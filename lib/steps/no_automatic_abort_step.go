package steps

type NoAutomaticAbort struct{}

func (step NoAutomaticAbort) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step NoAutomaticAbort) ShouldAutomaticallyAbortOnError() bool {
	return false
}
