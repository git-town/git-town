package steps

// NoAbortStep is a partial Step implementation used for composition
type NoAbortStep struct{}

// CreateAbortStep returns the abort step for this step.
func (step NoAbortStep) CreateAbortStep() Step {
	return NoOpStep{}
}
