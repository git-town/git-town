package steps

// SkipCurrentBranchSteps is a mock step to be used instead of
// running another list of steps.
// This is used when ignoring the remaining steps for a particular branch.
type SkipCurrentBranchSteps struct {
	NoExpectedError
	NoUndoStep
}

// Run executes this step.
func (step SkipCurrentBranchSteps) Run() error {
	return nil
}
