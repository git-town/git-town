package shared

// Opcode is an atomic operation that the Git Town interpreter can execute.
// Opcodes implement the command pattern (https://en.wikipedia.org/wiki/Command_pattern)
// and provide opcodes to continue and abort them.
// Undoing an opcode is done via the undo package.
type Opcode interface {
	// AbortProgram provides the opcodes to abort this Opcode when it encounters an error.
	AbortProgram() []Opcode

	// ContinueProgram provides the opcodes continue this opcode
	// after it encountered an error and the user has resolved the error.
	ContinueProgram() []Opcode

	// AutomaticUndoError provides the error message to display when this opcode
	// cause the command to automatically undo.
	AutomaticUndoError() error

	// Run executes this opcodes.
	Run(args RunArgs) error

	// ShouldUndoOnError indicates whether this opcode should
	// cause the command to automatically undo if it errors.
	// When true, automatically runs the abort and undo logic and leaves the user where they started.
	// When false, stops execution to let the user fix the issue and continue or manually undo.
	ShouldUndoOnError() bool

	// UndoProgram provides the opcodes to undo this operation.
	// All Git changes are automatically undone by the snapshot-based undo engine
	// and don't need to be undone here.
	// The undo program returned here is only for external changes
	// like updating proposals at the code hosting platform.
	UndoExternalChangesProgram() []Opcode
}
