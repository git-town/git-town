package steps

type NoUndoStepBeforeRun struct{}

func (step NoUndoStepBeforeRun) CreateUndoStepBeforeRun() Step {
	return NoOpStep{}
}
