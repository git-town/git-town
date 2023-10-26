package opcode_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/program"
	"github.com/git-town/git-town/v9/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestIfBranchHasUnmergedChanges(t *testing.T) {
	t.Parallel()

	t.Run("equal values", func(t *testing.T) {
		t.Parallel()
		one := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		two := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
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
		one := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		two := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
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
		one := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		two := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.ContinueRebase{},
			},
		}
		must.NotEq(t, one, two)
	})

	t.Run("different Branch", func(t *testing.T) {
		t.Parallel()
		one := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		two := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch-2"),
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		must.Eq(t, one, two)
	})

	t.Run("embedded in list", func(t *testing.T) {
		t.Parallel()
		one := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		two := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		list1 := []shared.Opcode{&one}
		list2 := []shared.Opcode{&two}
		must.Eq(t, list1, list2)
	})

	t.Run("embedded in program", func(t *testing.T) {
		t.Parallel()
		one := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		two := opcode.IfBranchHasUnmergedChanges{
			Branch: domain.NewLocalBranchName("branch"),
			WhenTrue: []shared.Opcode{
				&opcode.AbortMerge{},
			},
			WhenFalse: []shared.Opcode{
				&opcode.AbortRebase{},
			},
		}
		prog1 := program.Program{Opcodes: []shared.Opcode{&one}}
		prog2 := program.Program{Opcodes: []shared.Opcode{&two}}
		must.Eq(t, prog1, prog2)
	})
}
