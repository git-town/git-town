package execute

import (
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/validate"
)

// LoadRepoSnapshot loads the initial snapshot of the Git repo.
func LoadRepoSnapshot(args LoadRepoSnapshotArgs) (gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	var branchesSnapshot gitdomain.BranchesSnapshot
	var err error
	stashSize, err := args.Backend.StashSize()
	if err != nil {
		return branchesSnapshot, stashSize, false, err
	}
	if args.ValidateNoOpenChanges {
		err = validate.NoOpenChanges(args.RepoStatus.OpenChanges)
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
	}
	if args.Fetch {
		var remotes gitdomain.Remotes
		remotes, err := args.Backend.Remotes()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline.Bool() {
			err = args.Frontend.Fetch()
			if err != nil {
				return branchesSnapshot, stashSize, false, err
			}
		}
		// must always reload the snapshot here because we fetched updates from the remote
		branchesSnapshot, err = args.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
	}
	if branchesSnapshot.IsEmpty() {
		branchesSnapshot, err = args.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
	}
	return branchesSnapshot, stashSize, false, err
}

type LoadRepoSnapshotArgs struct {
	Backend               git.BackendCommands
	DialogTestInputs      components.TestInputs
	Fetch                 bool
	Frontend              git.FrontendCommands
	Repo                  OpenRepoResult
	RepoStatus            gitdomain.RepoStatus
	ValidateNoOpenChanges bool
}
