package steps

type PushBranchAfterCurrentBranchSteps struct{}

func (step PushBranchAfterCurrentBranchSteps) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step PushBranchAfterCurrentBranchSteps) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step PushBranchAfterCurrentBranchSteps) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step PushBranchAfterCurrentBranchSteps) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step PushBranchAfterCurrentBranchSteps) Run() error {
	return nil
}

func (step PushBranchAfterCurrentBranchSteps) ShouldAutomaticallyAbortOnError() bool {
	return false
}
