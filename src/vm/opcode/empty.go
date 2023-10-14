package opcode

import (
	"errors"
)

// undeclaredOpcodeMethods makes structs in this package satisfy the Opcode interface even if they don't declare all required methods.
type undeclaredOpcodeMethods struct{}

func (step *undeclaredOpcodeMethods) CreateAbortProgram() []Opcode {
	return []Opcode{}
}

func (step *undeclaredOpcodeMethods) CreateContinueProgram() []Opcode {
	return []Opcode{}
}

func (step *undeclaredOpcodeMethods) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (step *undeclaredOpcodeMethods) Run(_ RunArgs) error {
	return nil
}

func (step *undeclaredOpcodeMethods) ShouldAutomaticallyAbortOnError() bool {
	return false
}
