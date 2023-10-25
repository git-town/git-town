package opcode

import (
	"reflect"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// IfElse allows running different opcodes based on a condition evaluated at runtime.
type IfElse struct {
	Condition func(*git.BackendCommands, config.Lineage) (bool, error)
	WhenTrue  []shared.Opcode // the opcodes to execute if the given branch is empty
	WhenFalse []shared.Opcode // the opcodes to execute if the given branch is not empty
	undeclaredOpcodeMethods
}

// This method makes comparison via https://github.com/google/go-cmp work in unit tests.
func (self IfElse) Equal(other IfElse) bool {
	return reflect.DeepEqual(self.WhenFalse, other.WhenFalse) &&
		reflect.DeepEqual(self.WhenTrue, other.WhenTrue)
}

func (self *IfElse) Run(args shared.RunArgs) error {
	condition, err := self.Condition(&args.Runner.Backend, args.Lineage)
	if err != nil {
		return err
	}
	if condition {
		args.PrependOpcodes(self.WhenTrue...)
	} else {
		args.PrependOpcodes(self.WhenFalse...)
	}
	return nil
}
