package shared

// Opcode represents a dedicated CLI activity.
// Git Town commands consist of many Opcode instances.
// Steps implement the command pattern (https://en.wikipedia.org/wiki/Command_pattern)
// and can provide Steps to continue, abort, and undo them.
type Opcode interface {
	// CreateAbortProgram provides the abort step for this step.
	CreateAbortProgram() []Opcode

	// CreateContinueProgram provides the continue step for this step.
	CreateContinueProgram() []Opcode

	// CreateAutomaticAbortError provides the error message to display when this step
	// cause the command to automatically abort.
	CreateAutomaticAbortError() error

	// Run executes this step.
	Run(args RunArgs) error

	// ShouldAutomaticallyAbortOnError indicates whether this step should
	// cause the command to automatically abort if it errors.
	// When true, automatically runs the abort logic and leaves the user where they started.
	// When false, stops execution to let the user fix the issue and continue or manually abort.
	ShouldAutomaticallyAbortOnError() bool
}
