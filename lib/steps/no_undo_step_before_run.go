package steps

// NoUndoStepBeforeRun is a partial Step implementation used for composition
type NoUndoStepBeforeRun struct{}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step NoUndoStepBeforeRun) CreateUndoStepBeforeRun() Step {
	return NoOpStep{}
}
