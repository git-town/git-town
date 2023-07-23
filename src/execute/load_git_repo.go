package execute

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/validate"
)

func LoadGitRepo(pr *git.ProdRunner, args LoadGitArgs) (branchesSyncStatus git.BranchesSyncStatus, currentBranch string, exit bool, err error) { //nolint:nonamedreturns
	isRepo, rootDir := pr.Backend.IsRepositoryUncached()
	if !isRepo {
		err = errors.New("this is not a Git repository")
		return
	}
	if args.ValidateNoOpenChanges {
		var hasOpenChanges bool
		hasOpenChanges, err = pr.Backend.HasOpenChanges()
		if err != nil {
			return
		}
		if hasOpenChanges {
			err = fmt.Errorf("you have uncommitted changes. Did you mean to commit them before shipping?")
			return
		}
	}
	isOffline, err := pr.Config.IsOffline()
	if err != nil {
		return
	}
	if args.ValidateIsOnline && isOffline {
		err = errors.New("this command requires an active internet connection")
		return
	}
	if args.Fetch {
		var hasOrigin bool
		hasOrigin, err = pr.Backend.HasOrigin()
		if err != nil {
			return
		}
		if hasOrigin && !isOffline {
			err = pr.Frontend.Fetch()
			if err != nil {
				return
			}
		}
	}
	branchesSyncStatus, currentBranch, err = pr.Backend.BranchesSyncStatus()
	if err != nil {
		return
	}
	pr.Backend.CurrentBranchCache.Set(currentBranch)
	if args.ValidateIsConfigured {
		err = validate.IsConfigured(&pr.Backend)
		if err != nil {
			return
		}
	}
	currentDirectory, err := os.Getwd()
	if err != nil {
		err = errors.New("cannot determine the current directory")
		return
	}
	if currentDirectory != rootDir {
		err = pr.Frontend.NavigateToDir(rootDir)
		if err != nil {
			return
		}
	}
	if args.HandleUnfinishedState {
		exit, err = validate.HandleUnfinishedState(pr, nil)
	}
	return
}

type LoadGitArgs struct {
	Fetch                 bool
	ValidateIsOnline      bool
	HandleUnfinishedState bool
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}
