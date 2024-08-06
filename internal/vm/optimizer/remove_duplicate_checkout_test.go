package optimizer_test

import (
	"testing"

	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/vm/opcodes"
	"github.com/git-town/git-town/v14/internal/vm/optimizer"
	"github.com/git-town/git-town/v14/internal/vm/program"
	"github.com/shoenig/test/must"
)

func TestRemoveDuplicateCheckout(t *testing.T) {
	t.Parallel()

	t.Run("duplicate checkout opcodes", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-1")},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-2")},
			&opcodes.AbortRebase{},
		}
		have := optimizer.RemoveDuplicateCheckout(give)
		want := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-2")},
			&opcodes.AbortRebase{},
		}
		must.Eq(t, want, have)
	})

	t.Run("duplicate checkout opcodes mixed with end-of-branch opcodes", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-1")},
			&opcodes.EndOfBranchProgram{},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-2")},
			&opcodes.EndOfBranchProgram{},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-3")},
			&opcodes.AbortRebase{},
		}
		have := optimizer.RemoveDuplicateCheckout(give)
		want := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.EndOfBranchProgram{},
			&opcodes.EndOfBranchProgram{},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-3")},
			&opcodes.AbortRebase{},
		}
		must.Eq(t, want, have)
	})

	t.Run("a mix of Checkout and CheckoutIfExists opcodes", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-1")},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-2")},
		}
		have := optimizer.RemoveDuplicateCheckout(give)
		want := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-2")},
		}
		must.Eq(t, want, have)
	})

	t.Run("no duplicate checkout opcodes", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.AbortRebase{},
		}
		have := optimizer.RemoveDuplicateCheckout(give)
		want := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.AbortRebase{},
		}
		must.Eq(t, want, have)
	})
}
