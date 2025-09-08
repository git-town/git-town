package shared

import (
	"fmt"
	"strings"
)

// This file defines the interfaces that opcodes can implement.

// Opcode represents an opcode about which we know nothing except that it is an opcode.
type Opcode interface{}

// Runnable defines methods that an opcode needs to implement to execute subshell commands.
type Runnable interface {
	// Run executes this opcode.
	Run(args RunArgs) error
}

// Abortable allows an opcode that executes a Git command that can fail
// to define custom steps that abort that failing Git command.
type Abortable interface {
	Abort() []Opcode
}

// Abortable allows an opcode that executes a Git command that can fail
// to define custom steps that safely abort that Git command when it fails.
// By default, opcodes continue by executing their Run method again.
type Continuable interface {
	Continue() []Opcode
}

// AutoUndoable allows an opcode that exacutes a Git command that when it fails,
// it should fail the entire Git Town command, to specify the error that Git Town displays.
type AutoUndoable interface {
	AutomaticUndoError() error
}

// UndoExternalChanges allows an opcodes that performs external changes,
// for example at a forge, to undo them when there is an error.
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
