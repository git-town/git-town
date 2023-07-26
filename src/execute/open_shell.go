package execute

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v9/src/cache"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/statistics"
	"github.com/git-town/git-town/v9/src/subshell"
	"github.com/git-town/git-town/v9/src/validate"
)

func OpenShell(args LoadArgs) (prodRunner git.ProdRunner, rootDir string, isOffline bool, exit bool, err error) {
	var stats Statistics
	if args.Debug {
		stats = &statistics.CommandsRun{CommandsCount: 0}
	} else {
		stats = &statistics.None{}
	}
	backendRunner := subshell.BackendRunner{
		Dir:     nil,
		Stats:   stats,
		Verbose: args.Debug,
	}
	backendCommands := git.BackendCommands{
		BackendRunner:      backendRunner,
		Config:             nil, // NOTE: initializing to nil here to validate the Git version before running any Git commands, setting to the correct value after that is done
		CurrentBranchCache: &cache.String{},
		RemoteBranchCache:  &cache.Strings{},
		RemotesCache:       &cache.Strings{},
		RootDirCache:       &cache.String{},
	}
	err = validate.HasGitVersion(&backendCommands)
	if err != nil {
		return
	}
	config := git.RepoConfig{
		GitTown: config.NewGitTown(backendRunner),
		DryRun:  false, // to bootstrap this, DryRun always gets initialized as false and later enabled if needed
	}
	backendCommands.Config = &config
	prodRunner = git.ProdRunner{
		Config:  config,
		Backend: backendCommands,
		Frontend: git.FrontendCommands{
			FrontendRunner:         NewFrontendRunner(args.OmitBranchNames, args.DryRun, prodRunner.Backend.CurrentBranch, stats),
			SetCachedCurrentBranch: backendCommands.CurrentBranchCache.Set,
		},
		Stats: stats,
	}
	if args.DryRun {
		prodRunner.Config.DryRun = true
	}
	if args.ValidateGitRepo {
		rootDir = backendCommands.RootDirectory()
		if rootDir == "" {
			err = errors.New(messages.RepoOutside)
			return
		}
	}
	if args.HandleUnfinishedState {
		exit, err = validate.HandleUnfinishedState(&prodRunner, nil, rootDir)
		if err != nil || exit {
			return
		}
	}
	if args.ValidateNoOpenChanges {
		err = validate.NoOpenChanges(&prodRunner.Backend)
		if err != nil {
			return
		}
	}
	isOffline, err = config.IsOffline()
	if err != nil {
		return
	}
	if args.ValidateIsOnline && isOffline {
		err = errors.New(messages.OfflineNotAllowed)
		return
	}
	if args.Fetch {
		var hasOrigin bool
		hasOrigin, err = backendCommands.HasOrigin()
		if err != nil {
			return
		}
		if hasOrigin && !isOffline {
			err = prodRunner.Frontend.Fetch()
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
	if args.ValidateGitRepo && currentDirectory != rootDir {
		err = prodRunner.Frontend.NavigateToDir(rootDir)
		if err != nil {
			return
		}
	}
	return
}

type LoadArgs struct {
	Debug                 bool
	DryRun                bool
	Fetch                 bool
	HandleUnfinishedState bool
	OmitBranchNames       bool
	ValidateGitRepo       bool
	ValidateIsOnline      bool
	ValidateNoOpenChanges bool
}

// NewFrontendRunner provides a FrontendRunner instance that behaves according to the given configuration.
func NewFrontendRunner(omitBranchNames, dryRun bool, getCurrentBranch subshell.GetCurrentBranchFunc, stats Statistics) git.FrontendRunner {
	if dryRun {
		return &subshell.FrontendDryRunner{
			GetCurrentBranch: getCurrentBranch,
			OmitBranchNames:  omitBranchNames,
			Stats:            stats,
		}
	}
	return &subshell.FrontendRunner{
		GetCurrentBranch: getCurrentBranch,
		OmitBranchNames:  omitBranchNames,
		Stats:            stats,
	}
}

type Statistics interface {
	RegisterRun()
	PrintAnalysis()
}
