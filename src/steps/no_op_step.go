package steps

// NoOpStep does nothing.
// It is used for steps that have no undo or abort steps.
type NoOpStep struct{}

// CreateAbortStep returns the abort step for this step.
func (step *NoOpStep) CreateAbortStep() Step {
	return &NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step *NoOpStep) CreateContinueStep() Step {
	return &NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step *NoOpStep) CreateUndoStep() Step {
	return &NoOpStep{}
}

// GetAutomaticAbortErrorMessage returns the error message to display when this step
// cause the command to automatically abort.
func (step *NoOpStep) GetAutomaticAbortErrorMessage() string {
	return ""
}

// Run executes this step.
func (step *NoOpStep) Run() error {
	return nil
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *NoOpStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
