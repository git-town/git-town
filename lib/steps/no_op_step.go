package steps

// NoOpStep does nothing.
// It is used for steps that have no undo or abort steps.
type NoOpStep struct {
	NoAutomaticAbortOnError
}

// CreateAbortStep returns the abort step for this step.
func (step NoOpStep) CreateAbortStep() Step {
	return step
}

// CreateContinueStep returns the continue step for this step.
func (step NoOpStep) CreateContinueStep() Step {
	return step
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step NoOpStep) CreateUndoStepBeforeRun() Step {
	return step
}

// CreateUndoStepAfterRun returns the undo step for this step after it is run.
func (step NoOpStep) CreateUndoStepAfterRun() Step {
	return step
}

// Run executes this step.
func (step NoOpStep) Run() error {
	return nil
}
