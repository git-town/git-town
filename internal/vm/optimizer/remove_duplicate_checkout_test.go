package optimizer_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/opcodes"
	"github.com/git-town/git-town/v17/internal/vm/optimizer"
	"github.com/git-town/git-town/v17/internal/vm/program"
	"github.com/shoenig/test/must"
)

func TestRemoveDuplicateCheckout(t *testing.T) {
	t.Parallel()

	t.Run("duplicate checkout opcodes", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-1")},
			&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-2")},
			&opcodes.RebaseAbort{},
		}
		have := optimizer.RemoveDuplicateCheckout(give)
		want := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-2")},
			&opcodes.RebaseAbort{},
		}
		must.Eq(t, want, have)
	})

	t.Run("duplicate checkout opcodes mixed with end-of-branch opcodes", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-1")},
			&opcodes.ProgramEndOfBranch{},
			&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-2")},
			&opcodes.ProgramEndOfBranch{},
			&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-3")},
			&opcodes.RebaseAbort{},
		}
		have := optimizer.RemoveDuplicateCheckout(give)
		want := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.ProgramEndOfBranch{},
			&opcodes.ProgramEndOfBranch{},
			&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-3")},
			&opcodes.RebaseAbort{},
		}
		must.Eq(t, want, have)
	})

	t.Run("a mix of Checkout and CheckoutIfExists opcodes", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-1")},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-2")},
		}
		have := optimizer.RemoveDuplicateCheckout(give)
		want := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-2")},
		}
		must.Eq(t, want, have)
	})

	t.Run("no duplicate checkout opcodes", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.RebaseAbort{},
		}
		have := optimizer.RemoveDuplicateCheckout(give)
		want := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.RebaseAbort{},
		}
		must.Eq(t, want, have)
	})
}
