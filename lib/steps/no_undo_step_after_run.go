package steps

// NoUndoStepAfterRun is a partial Step implementation used for composition
type NoUndoStepAfterRun struct{}

// CreateUndoStepAfterRun returns the undo step for this step after it is run.
func (step NoUndoStepAfterRun) CreateUndoStepAfterRun() Step {
	return NoOpStep{}
}
