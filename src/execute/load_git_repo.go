package execute

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/validate"
)

func LoadGitRepo(pr *git.ProdRunner, args LoadGitArgs) (rootDir string, isOffline bool, exit bool, err error) { //nolint:nonamedreturns
	rootDir = pr.Backend.RootDirectory()
	if rootDir == "" {
		err = errors.New(messages.RepoOutside)
		return
	}
	if args.HandleUnfinishedState {
		exit, err = validate.HandleUnfinishedState(pr, nil, rootDir)
		if err != nil || exit {
			return
		}
	}
	if args.ValidateNoOpenChanges {
		err = validate.NoOpenChanges(pr.Backend)
		if err != nil {
			return
		}
	}
	isOffline, err = pr.Config.IsOffline()
	if err != nil {
		return
	}
	if args.ValidateIsOnline && isOffline {
		err = errors.New(messages.OfflineNotAllowed)
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
	currentDirectory, err := os.Getwd()
	if err != nil {
		err = errors.New(messages.DirCurrentProblem)
		return
	}
	if currentDirectory != rootDir {
		err = pr.Frontend.NavigateToDir(rootDir)
		if err != nil {
			return
		}
	}
	return
}

type LoadGitArgs struct {
	Fetch                 bool
	ValidateIsOnline      bool
	HandleUnfinishedState bool
	ValidateNoOpenChanges bool
}
