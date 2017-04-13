package steps

// NoContinueStep is a partial Step implementation used for composition
type NoContinueStep struct{}

// CreateContinueStep returns the continue step for this step.
func (step NoContinueStep) CreateContinueStep() Step {
	return NoOpStep{}
}
