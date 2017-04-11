package steps

// NoOpStep does nothing.
// It is used for steps that have no undo or abort steps.
type NoOpStep struct{}

// CreateAbortStep returns the abort step for this step.
func (step NoOpStep) CreateAbortStep() Step {
	return step
}

// CreateContinueStep returns the continue step for this step.
func (step NoOpStep) CreateContinueStep() Step {
	return step
}

// CreateUndoStep returns the undo step for this step.
func (step NoOpStep) CreateUndoStep() Step {
	return step
}

// Run executes this step.
func (step NoOpStep) Run() error {
	return nil
}
