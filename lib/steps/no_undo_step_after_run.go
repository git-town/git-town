package steps

type NoUndoStepAfterRun struct{}

func (step NoUndoStepAfterRun) CreateUndoStepAfterRun() Step {
	return NoOpStep{}
}
