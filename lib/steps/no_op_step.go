package steps

type NoOpStep struct {
	NoAutomaticAbortOnError
}

func (step NoOpStep) CreateAbortStep() Step {
	return step
}

func (step NoOpStep) CreateContinueStep() Step {
	return step
}

func (step NoOpStep) CreateUndoStepBeforeRun() Step {
	return step
}

func (step NoOpStep) CreateUndoStepAfterRun() Step {
	return step
}

func (step NoOpStep) Run() error {
	return nil
}
