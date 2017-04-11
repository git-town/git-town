package steps

type PushBranchAfterCurrentBranchSteps struct {
	NoAutomaticAbortOnError
	NoUndoStep
}

func (step PushBranchAfterCurrentBranchSteps) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step PushBranchAfterCurrentBranchSteps) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step PushBranchAfterCurrentBranchSteps) Run() error {
	return nil
}
