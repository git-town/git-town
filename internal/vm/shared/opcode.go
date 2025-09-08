package shared

import (
	"fmt"
	"strings"
)

// The common methods for all opcodes.
type Opcode interface{}

// Runnable allows an opcode to execute subshell commands.
type Runnable interface {
	// Run executes this opcodes.
	Run(args RunArgs) error
}

// Recoverable defines methods for opcodes that can encounter conflicts.
type Recoverable interface {
	// Abort provides the opcodes to abort this Opcode when it encounters an error.
	Abort() []Opcode

	// Continue provides the opcodes continue this opcode
	// after it encountered an error and the user has resolved the error.
	Continue() []Opcode
}

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
