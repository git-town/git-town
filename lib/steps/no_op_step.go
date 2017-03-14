package steps

type NoOpStep struct{}

func (step NoOpStep) CreateAbortStep() Step {
	return step
}

func (step NoOpStep) CreateContinueStep() Step {
	return step
}

func (step NoOpStep) CreateUndoStep() Step {
	return step
}

func (step NoOpStep) Run() error {
	return nil
}
