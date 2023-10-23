package opcode_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestIfElse(t *testing.T) {
	t.Parallel()
	one := opcode.IfElse{
		Condition: func(bc *git.BackendCommands, l config.Lineage) (bool, error) { return true, nil },
		WhenFalse: []shared.Opcode{},
		WhenTrue:  []shared.Opcode{},
	}
	two := opcode.IfElse{
		Condition: func(bc *git.BackendCommands, l config.Lineage) (bool, error) { return true, nil },
		WhenFalse: []shared.Opcode{},
		WhenTrue:  []shared.Opcode{},
	}
	must.Eq(t, one, two)
}
