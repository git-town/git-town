package steps

// SkipCurrentBranchSteps is a mock step to be used instead of
// running another list of steps.
// This is used when ignoring the remaining steps for a particular branch.
type SkipCurrentBranchSteps struct{}

// CreateAbortStep returns the abort step for this step.
func (step SkipCurrentBranchSteps) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step SkipCurrentBranchSteps) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step SkipCurrentBranchSteps) CreateUndoStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step SkipCurrentBranchSteps) Run() error {
	return nil
}
