package shared_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	"github.com/shoenig/test/must"
)

func TestBranchesInOpcode(t *testing.T) {
	t.Parallel()

	t.Run("BranchName", func(t *testing.T) {
		t.Parallel()
		opcode := opcodes.MergeIntoCurrentBranch{
			BranchToMerge: "branch",
		}
		have := shared.BranchesInOpcode(&opcode)
		want := []gitdomain.BranchName{
			"branch",
		}
		must.Eq(t, want, have)
	})

	t.Run("LocalBranchName", func(t *testing.T) {
		t.Parallel()
		opcode := opcodes.LineageParentSet{
			Branch: "branch",
			Parent: "parent",
		}
		have := shared.BranchesInOpcode(&opcode)
		want := []gitdomain.BranchName{
			"branch",
			"parent",
		}
		must.Eq(t, want, have)
	})

	t.Run("LocalBranchNames", func(t *testing.T) {
		t.Parallel()
		opcode := opcodes.LineageParentSetFirstExisting{
			Branch:    "branch",
			Ancestors: gitdomain.NewLocalBranchNames("ancestor-1", "ancestor-2"),
		}
		have := shared.BranchesInOpcode(&opcode)
		want := []gitdomain.BranchName{
			"ancestor-1",
			"ancestor-2",
			"branch",
		}
		must.Eq(t, want, have)
	})

	t.Run("RemoteBranchName", func(t *testing.T) {
		t.Parallel()
		opcode := opcodes.BranchTrackingDelete{
			Branch: "origin/branch",
		}
		have := shared.BranchesInOpcode(&opcode)
		want := []gitdomain.BranchName{
			"origin/branch",
		}
		must.Eq(t, want, have)
	})
}
