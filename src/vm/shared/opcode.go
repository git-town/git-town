package shared

// Opcode is an atomic operation that the Git Town interpreter can execute.
// Opcodes implement the command pattern (https://en.wikipedia.org/wiki/Command_pattern)
// and provide opcodes to continue and abort them.
// Undoing an opcode is done via the undo package.
type Opcode interface {
	// CreateAbortProgram provides the opcodes to abort this Opcode when it encounters an error.
	CreateAbortProgram() []Opcode

	// CreateContinueProgram provides the opcodes continue this opcode
	// after it encountered an error and the user has resolved the error.
	CreateContinueProgram() []Opcode

	// CreateAutomaticUndoError provides the error message to display when this opcode
	// cause the command to automatically undo.
	CreateAutomaticUndoError() error

	// Run executes this opcodes.
	Run(args RunArgs) error

	// ShouldAutomaticallyUndoOnError indicates whether this opcode should
	// cause the command to automatically undo if it errors.
	// When true, automatically runs the abort and undo logic and leaves the user where they started.
	// When false, stops execution to let the user fix the issue and continue or manually undo.
	ShouldAutomaticallyUndoOnError() bool
}
