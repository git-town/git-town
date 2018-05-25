package steps

// Step represents a dedicated activity within a Git Town command.
// Git Town commands are comprised of a number of steps that need to be executed.
type Step interface {
	CreateAbortStep() Step
	CreateContinueStep() Step
	CreateUndoStepBeforeRun() Step
	CreateUndoStepAfterRun() Step
	GetAutomaticAbortErrorMessage() string
	Run() error
	ShouldAutomaticallyAbortOnError() bool
}
