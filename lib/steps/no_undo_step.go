package steps

// NoUndoStep is a partial Step implementation used for composition
type NoUndoStep struct {
	NoUndoStepAfterRun
	NoUndoStepBeforeRun
}
