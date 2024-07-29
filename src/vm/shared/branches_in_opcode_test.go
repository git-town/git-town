package shared_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestBranchesInOpcode(t *testing.T) {
	t.Parallel()

	t.Run("LocalBranchName", func(t *testing.T) {
		t.Parallel()
		opcode := opcodes.ChangeParent{
			Branch: gitdomain.NewLocalBranchName("branch"),
			Parent: gitdomain.NewLocalBranchName("parent"),
		}
		have := shared.BranchesInOpcode(&opcode)
		want := []gitdomain.BranchName{
			gitdomain.NewBranchName("branch"),
			gitdomain.NewBranchName("parent"),
		}
		must.Eq(t, want, have)
	})

	t.Run("RemoteBranchName", func(t *testing.T) {
		t.Parallel()
		opcode := opcodes.DeleteTrackingBranch{
			Branch: gitdomain.NewRemoteBranchName("origin/branch"),
		}
		have := shared.BranchesInOpcode(&opcode)
		want := []gitdomain.BranchName{
			gitdomain.NewBranchName("origin/branch"),
		}
		must.Eq(t, want, have)
	})

	t.Run("BranchName", func(t *testing.T) {
		t.Parallel()
		opcode := opcodes.Merge{
			Branch: gitdomain.NewBranchName("branch"),
		}
		have := shared.BranchesInOpcode(&opcode)
		want := []gitdomain.BranchName{
			gitdomain.NewBranchName("branch"),
		}
		must.Eq(t, want, have)
	})
}
