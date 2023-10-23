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

	t.Run("equal values", func(t *testing.T) {
		t.Parallel()
		one := opcode.IfElse{
			Condition: func(bc *git.BackendCommands, l config.Lineage) (bool, error) { return true, nil },
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		two := opcode.IfElse{
			Condition: func(bc *git.BackendCommands, l config.Lineage) (bool, error) { return true, nil },
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		must.Eq(t, one, two)
	})

	t.Run("different WhenTrue values", func(t *testing.T) {
		t.Parallel()
		one := opcode.IfElse{
			Condition: func(bc *git.BackendCommands, l config.Lineage) (bool, error) { return true, nil },
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		two := opcode.IfElse{
			Condition: func(bc *git.BackendCommands, l config.Lineage) (bool, error) { return true, nil },
			WhenTrue: []shared.Opcode{
				&opcode.ContinueMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		must.NotEq(t, one, two)
	})

	t.Run("different WhenFalse values", func(t *testing.T) {
		t.Parallel()
		one := opcode.IfElse{
			Condition: func(bc *git.BackendCommands, l config.Lineage) (bool, error) { return true, nil },
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		two := opcode.IfElse{
			Condition: func(bc *git.BackendCommands, l config.Lineage) (bool, error) { return true, nil },
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.ContinueRebase{},
			},
		}
		must.NotEq(t, one, two)
	})

	t.Run("different condition function", func(t *testing.T) {
		t.Parallel()
		one := opcode.IfElse{
			Condition: func(bc *git.BackendCommands, l config.Lineage) (bool, error) { return true, nil },
			WhenTrue: []shared.Opcode{
				&opcode.AbortRebase{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortMerge{},
			},
		}
		two := opcode.IfElse{
			Condition: func(bc *git.BackendCommands, l config.Lineage) (bool, error) { return false, nil },
			WhenTrue: []shared.Opcode{
				&opcode.AbortRebase{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortMerge{},
			},
		}
		must.Eq(t, one, two)
	})
}
