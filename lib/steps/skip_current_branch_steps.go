package steps

type SkipCurrentBranchSteps struct {
	NoAutomaticAbortOnError
	NoUndoStep
}

func (step SkipCurrentBranchSteps) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step SkipCurrentBranchSteps) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step SkipCurrentBranchSteps) Run() error {
	return nil
}
