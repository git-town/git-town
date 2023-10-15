package shared

// Opcode is an atomic operation that the Git Town interpreter can execute.
// Opcodes implement the command pattern (https://en.wikipedia.org/wiki/Command_pattern)
// and provide opcodes to continue and abort them.
// Undoing an opcode is done via the undo package.
type Opcode interface {
	// CreateAbortProgram provides the abort opcodes for this opcode.
	CreateAbortProgram() []Opcode

	// CreateContinueProgram provides the continue opcodes for this opcode.
	CreateContinueProgram() []Opcode

	// CreateAutomaticAbortError provides the error message to display when this opcode
	// cause the command to automatically abort.
	CreateAutomaticAbortError() error

	// Run executes this opcode.
	Run(args RunArgs) error

	// ShouldAutomaticallyAbortOnError indicates whether this opcode should
	// cause the command to automatically abort if it errors.
	// When true, automatically runs the abort logic and leaves the user where they started.
	// When false, stops execution to let the user fix the issue and continue or manually abort.
	ShouldAutomaticallyAbortOnError() bool
}
