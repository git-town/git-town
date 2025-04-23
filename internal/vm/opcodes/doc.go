// Package opcodes defines the individual operations that the Git Town VM can execute.
// All opcodes implement the shared.Opcode interface.
package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v19/internal/gohacks"
	"github.com/git-town/git-town/v19/internal/vm/shared"
)

// undeclaredOpcodeMethods can be added to structs in this package to satisfy the shared.Opcode interface even if they don't declare all required methods.
type undeclaredOpcodeMethods struct{}

func (self *undeclaredOpcodeMethods) AbortProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (self *undeclaredOpcodeMethods) AutomaticUndoError() error {
	return errors.New("")
}

func (self *undeclaredOpcodeMethods) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (self *undeclaredOpcodeMethods) Run(_ shared.RunArgs) error {
	return nil
}

func (self *undeclaredOpcodeMethods) ShouldUndoOnError() bool {
	return false
}

func (self *undeclaredOpcodeMethods) UndoExternalChangesProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func Lookup(opcodeType string) shared.Opcode { //nolint:ireturn
	for _, opcode := range All() {
		if gohacks.TypeName(opcode) == opcodeType {
			return opcode
		}
	}
	return nil
}
