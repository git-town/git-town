package opcode

import (
	"errors"
)

// BaseOpcode does nothing.
// It is used for steps that have no undo or abort steps.
type BaseOpcode struct{}

func (step *BaseOpcode) CreateAbortProgram() []Opcode {
	return []Opcode{}
}

func (step *BaseOpcode) CreateContinueProgram() []Opcode {
	return []Opcode{}
}

func (step *BaseOpcode) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (step *BaseOpcode) Run(_ RunArgs) error {
	return nil
}

func (step *BaseOpcode) ShouldAutomaticallyAbortOnError() bool {
	return false
}
