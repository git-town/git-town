package opcode

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// IfElse allows running different steps based on a condition evaluated at runtime.
type IfElse struct {
	Condition func(*git.BackendCommands, config.Lineage) (bool, error)
	WhenTrue  []shared.Opcode // the steps to execute if the given branch is empty
	WhenFalse []shared.Opcode // the steps to execute if the given branch is not empty
	undeclaredOpcodeMethods
}

func (step *IfElse) Run(args shared.RunArgs) error {
	condition, err := step.Condition(&args.Runner.Backend, args.Lineage)
	if err != nil {
		return err
	}
	if condition {
		args.PrependOpcodes(step.WhenTrue...)
	} else {
		args.PrependOpcodes(step.WhenFalse...)
	}
	return nil
}
