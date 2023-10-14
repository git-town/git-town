package opcode

import (
	"errors"
)

// Empty does nothing.
// It is used for steps that have no undo or abort steps.
type Empty struct{}

func (step *Empty) CreateAbortProgram() []Opcode {
	return []Opcode{}
}

func (step *Empty) CreateContinueProgram() []Opcode {
	return []Opcode{}
}

func (step *Empty) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (step *Empty) Run(_ RunArgs) error {
	return nil
}

func (step *Empty) ShouldAutomaticallyAbortOnError() bool {
	return false
}
