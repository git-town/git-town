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

func OpenRepo(args OpenRepoArgs) (result OpenRepoResult, err error) {
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
		CurrentBranchCache: &cache.LocalBranch{},
		RemoteBranchCache:  &cache.RemoteBranch{},
		RemotesCache:       &cache.Remotes{},
	}
	majorVersion, minorVersion, err := backendCommands.Version()
	if err != nil {
		return result, err
	}
	err = validate.HasGitVersion(majorVersion, minorVersion)
	if err != nil {
		return
	}
	repoConfig := git.RepoConfig{
		GitTown: config.NewGitTown(backendRunner),
		DryRun:  false, // to bootstrap this, DryRun always gets initialized as false and later enabled if needed
	}
	backendCommands.Config = &repoConfig
	prodRunner := git.ProdRunner{
		Config:  repoConfig,
		Backend: backendCommands,
		Frontend: git.FrontendCommands{
			FrontendRunner:         NewFrontendRunner(args.OmitBranchNames, args.DryRun, backendCommands.CurrentBranch, stats),
			SetCachedCurrentBranch: backendCommands.CurrentBranchCache.Set,
		},
		Stats: stats,
	}
	if args.DryRun {
		prodRunner.Config.DryRun = true
	}
	rootDir := backendCommands.RootDirectory()
	if args.ValidateGitRepo {
		if rootDir == "" {
			err = errors.New(messages.RepoOutside)
			return
		}
	}
	isOffline, err := repoConfig.IsOffline()
	if err != nil {
		return
	}
	if args.ValidateIsOnline && isOffline {
		err = errors.New(messages.OfflineNotAllowed)
		return
	}
	if args.ValidateGitRepo {
		var currentDirectory string
		currentDirectory, err = os.Getwd()
		if err != nil {
			err = errors.New(messages.DirCurrentProblem)
			return
		}
		if currentDirectory != rootDir {
			err = prodRunner.Frontend.NavigateToDir(rootDir)
		}
	}
	return OpenRepoResult{
		Runner:    prodRunner,
		RootDir:   rootDir,
		IsOffline: isOffline,
	}, err
}

type OpenRepoArgs struct {
	Debug            bool
	DryRun           bool
	OmitBranchNames  bool
	ValidateGitRepo  bool
	ValidateIsOnline bool
}

type OpenRepoResult struct {
	Runner    git.ProdRunner
	RootDir   string
	IsOffline bool
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
