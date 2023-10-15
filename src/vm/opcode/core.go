// Package opcode defines the individual operations that the Git Town VM can execute.
// All opcodes implement the shared.Opcode interface.
package opcode

import (
	"errors"

	"github.com/git-town/git-town/v9/src/vm/shared"
)

// undeclaredOpcodeMethods can be added to structs in this package to satisfy the shared.Opcode interface even if they don't declare all required methods.
type undeclaredOpcodeMethods struct{}

func (op *undeclaredOpcodeMethods) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (op *undeclaredOpcodeMethods) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{}
}

func (op *undeclaredOpcodeMethods) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (op *undeclaredOpcodeMethods) Run(_ shared.RunArgs) error {
	return nil
}

func (op *undeclaredOpcodeMethods) ShouldAutomaticallyAbortOnError() bool {
	return false
}
