package shared

import (
	"fmt"
	"strings"
)

// Opcode represents an opcode about which we know nothing except that it is an opcode.
type Opcode interface{}

// Runnable defines methods that an opcode needs to implement to execute subshell commands.
type Runnable interface {
	// Run executes this opcodes.
	Run(args RunArgs) error
}

// Abortable allows an opcode that executes a Git command that can fail
// to define custom steps that safely abort that Git command when it fails.
type Abortable interface {
	// Abort provides the opcodes to abort this Opcode when it encounters an error.
	Abort() []Opcode
}

// Abortable allows an opcode that executes a Git command that can fail
// to define custom steps that safely abort that Git command when it fails.
// By default, opcodes continue by executing their Run method again.
type Continuable interface {
	// Continue provides the opcodes continue this opcode
	// after it encountered an error and the user has resolved the error.
	Continue() []Opcode
}

// AutoUndoable allows an opcode that exacutes a Git command that can fail
// to specify that it should
type AutoUndoable interface {
	// AutomaticUndoError provides the error message to display when this opcode
	// cause the command to automatically undo.
	AutomaticUndoError() error
}

type ExternalEffects interface {
	// UndoExternalChanges provides the opcodes to undo this operation.
	// All Git changes are automatically undone by the snapshot-based undo engine
	// and don't need to be undone here.
	// The undo program returned here is only for external changes
	// like updating proposals at the forge.
	UndoExternalChanges() []Opcode
}

func RenderOpcodes(opcodes []Opcode, indent string) string {
	sb := strings.Builder{}
	if len(opcodes) == 0 {
		sb.WriteString("(empty program)\n")
	} else {
		sb.WriteString("Program:\n")
		for o, opcode := range opcodes {
			sb.WriteString(fmt.Sprintf("%s%d: %#v\n", indent, o+1, opcode))
		}
	}
	return sb.String()
}
