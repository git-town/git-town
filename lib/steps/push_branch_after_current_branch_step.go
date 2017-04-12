package steps

// PushBranchAfterCurrentBranchSteps is a mock step that is used in the undo process
// to push the branch after other steps have been undone.
type PushBranchAfterCurrentBranchSteps struct {
	NoAutomaticAbortOnError
	NoUndoStep
}

// CreateAbortStep returns the abort step for this step.
func (step PushBranchAfterCurrentBranchSteps) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step PushBranchAfterCurrentBranchSteps) CreateContinueStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step PushBranchAfterCurrentBranchSteps) Run() error {
	return nil
}
