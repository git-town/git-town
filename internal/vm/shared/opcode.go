package shared

// This file defines the interfaces that opcodes can implement.

// Opcode represents an opcode about which we know nothing except that it is an opcode.
// Opcode is an atomic operation that the Git Town interpreter can execute.
// Opcodes implement the command pattern (https://en.wikipedia.org/wiki/Command_pattern)
// and provide opcodes to continue and abort them.
// Undoing an opcode is done via the undo package.
type Opcode any

// Runnable marks an opcode that can execute subshell commands.
type Runnable interface {
	// Run executes this opcode.
	Run(args RunArgs) error
}

// Abortable marks an opcode that can provide custom steps
// to abort a failed Git command.
type Abortable interface {
	Abort() []Opcode
}

// Continuable marks an opcode that can provide custom steps
// to safely continue after a failed Git command.
// By default, opcodes retry by running their Run method again.
type Continuable interface {
	Continue() []Opcode
}

// AutoUndoable marks an opcode that should fail the entire Git Town command
// when its Git command fails. It provides the error message Git Town displays.
type AutoUndoable interface {
	AutomaticUndoError() error
}

// ExternalEffects marks an opcode that performs side effects outside of Git,
// such as changes on a code hosting service, and can undo them if needed.
type ExternalEffects interface {
	UndoExternalChanges() []Opcode
}
