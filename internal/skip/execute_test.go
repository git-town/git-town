package skip_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/skip"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	"github.com/shoenig/test/must"
)

func TestRemoveOpcodesForCurrentBranch(t *testing.T) {
	t.Parallel()

	t.Run("program contains multiple branches", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.Checkout{Branch: "branch-1"},
			&opcodes.PullCurrentBranch{},
			&opcodes.ProgramEndOfBranch{},
			&opcodes.Checkout{Branch: "branch-2"},
			&opcodes.PullCurrentBranch{},
			&opcodes.ProgramEndOfBranch{},
			&opcodes.Checkout{Branch: "branch-3"},
			&opcodes.PullCurrentBranch{},
			&opcodes.ProgramEndOfBranch{},
		}
		have := skip.RemoveOpcodesForCurrentBranch(give)
		want := program.Program{
			&opcodes.Checkout{Branch: "branch-2"},
			&opcodes.PullCurrentBranch{},
			&opcodes.ProgramEndOfBranch{},
			&opcodes.Checkout{Branch: "branch-3"},
			&opcodes.PullCurrentBranch{},
			&opcodes.ProgramEndOfBranch{},
		}
		must.Eq(t, want.String(), have.String())
	})

	t.Run("program contains no end of branch markers", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.Checkout{Branch: "branch-1"},
			&opcodes.PullCurrentBranch{},
			&opcodes.Checkout{Branch: "branch-2"},
			&opcodes.PullCurrentBranch{},
			&opcodes.Checkout{Branch: "branch-3"},
			&opcodes.PullCurrentBranch{},
		}
		have := skip.RemoveOpcodesForCurrentBranch(give)
		want := program.Program{}
		must.Eq(t, want.String(), have.String())
	})
}
