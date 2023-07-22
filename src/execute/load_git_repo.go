package execute

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/validate"
)

func LoadGitRepo(pr *git.ProdRunner, args LoadGitArgs) (branchesSyncStatus git.BranchesSyncStatus, currentBranch string, exit bool, err error) { //nolint:nonamedreturns
	fc := failure.Collector{}
	if args.Fetch {
		hasOrigin := fc.Bool(pr.Backend.HasOrigin())
		isOffline := fc.Bool(pr.Config.IsOffline())
		if fc.Err == nil && hasOrigin && !isOffline {
			fc.Check(pr.Frontend.Fetch())
		}
	}
	branchesSyncStatus, currentBranch, err = pr.Backend.BranchesSyncStatus()
	if err != nil {
		return branchesSyncStatus, currentBranch, false, errors.New("this is not a Git repository")
	}
	currentDirectory, err := os.Getwd()
	if err != nil {
		return branchesSyncStatus, currentBranch, false, fmt.Errorf("cannot get current working directory: %w", err)
	}
	gitRootDirectory, err := pr.Backend.RootDirectory()
	if err != nil {
		return branchesSyncStatus, currentBranch, false, err
	}
	if currentDirectory != gitRootDirectory {
		err = pr.Frontend.NavigateToDir(gitRootDirectory)
	}
	if args.ValidateIsOnline {
		fc.Check(validate.IsOnline(&pr.Config))
	}
	if args.HandleUnfinishedState {
		exit = fc.Bool(validate.HandleUnfinishedState(pr, nil))
	}
	return branchesSyncStatus, currentBranch, exit, fc.Err
}

type LoadGitArgs struct {
	Fetch                 bool
	ValidateIsOnline      bool
	HandleUnfinishedState bool
}
