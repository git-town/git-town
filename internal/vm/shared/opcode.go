package shared

import (
	"fmt"
	"strings"
)

// This file defines the interfaces that opcodes can implement.

// Opcode represents an opcode about which we know nothing except that it is an opcode.
type Opcode interface{}

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
