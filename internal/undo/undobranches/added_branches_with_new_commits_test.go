package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/undo/undobranches"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestAddedBranchesWithNewCommits(t *testing.T) {
	t.Parallel()

	mainBranch := gitdomain.BranchInfo{
		Local:      Some(gitdomain.BranchData{Name: "main", SHA: "111111"}),
		RemoteName: None[gitdomain.RemoteBranchName](),
		RemoteSHA:  None[gitdomain.SHA](),
		SyncStatus: gitdomain.SyncStatusLocalOnly,
	}
	begin := gitdomain.BranchesSnapshot{
		Active:   gitdomain.LocalBranchNameOpt("main"),
		Branches: gitdomain.BranchInfos{mainBranch},
	}
	end := gitdomain.BranchesSnapshot{
		Active: gitdomain.LocalBranchNameOpt("feature"),
		Branches: gitdomain.BranchInfos{
			mainBranch,
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "feature", SHA: "111111"}),
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
			},
		},
	}
	touched := []gitdomain.BranchName{"main", "feature"}

	t.Run("added branch is unchanged", func(t *testing.T) {
		t.Parallel()
		have := undobranches.AddedBranchesWithNewCommits(begin, end, touched, end.Branches)
		must.Eq(t, gitdomain.LocalBranchNames{}, have)
	})

	t.Run("added branch no longer exists", func(t *testing.T) {
		t.Parallel()
		current := gitdomain.BranchInfos{mainBranch}
		have := undobranches.AddedBranchesWithNewCommits(begin, end, touched, current)
		must.Eq(t, gitdomain.LocalBranchNames{}, have)
	})

	t.Run("added branch received new commits", func(t *testing.T) {
		t.Parallel()
		current := gitdomain.BranchInfos{
			mainBranch,
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "feature", SHA: "222222"}),
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
			},
		}
		have := undobranches.AddedBranchesWithNewCommits(begin, end, touched, current)
		want := gitdomain.NewLocalBranchNames("feature")
		must.Eq(t, want, have)
	})

	t.Run("added branch was not touched by the command", func(t *testing.T) {
		t.Parallel()
		current := gitdomain.BranchInfos{
			mainBranch,
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "feature", SHA: "222222"}),
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
			},
		}
		have := undobranches.AddedBranchesWithNewCommits(begin, end, []gitdomain.BranchName{"main"}, current)
		must.Eq(t, gitdomain.LocalBranchNames{}, have)
	})
}
