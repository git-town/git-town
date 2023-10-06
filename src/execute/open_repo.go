package execute

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/gohacks/cache"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/subshell"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/git-town/git-town/v9/src/validate"
)

func OpenRepo(args OpenRepoArgs) (*OpenRepoResult, error) {
	commandsCounter := gohacks.Counter{}
	backendRunner := subshell.BackendRunner{
		Dir:             nil,
		CommandsCounter: &commandsCounter,
		Verbose:         args.Debug,
	}
	backendCommands := git.BackendCommands{
		BackendRunner:      backendRunner,
		Config:             nil, // initializing to nil here to validate the Git version before running any Git commands, setting to the correct value after that is done
		CurrentBranchCache: &cache.LocalBranch{},
		RemotesCache:       &cache.Remotes{},
	}
	majorVersion, minorVersion, err := backendCommands.Version()
	if err != nil {
		return nil, err
	}
	err = validate.HasGitVersion(majorVersion, minorVersion)
	if err != nil {
		return nil, err
	}
	currentDirectory, err := os.Getwd()
	if err != nil {
		err = errors.New(messages.DirCurrentProblem)
		return nil, err
	}
	configSnapshot := undo.ConfigSnapshot{
		Cwd:       currentDirectory,
		GitConfig: config.LoadGitConfig(backendRunner),
	}
	repoConfig := git.RepoConfig{
		GitTown: config.NewGitTown(configSnapshot.GitConfig.Clone(), backendRunner),
		DryRun:  false, // to bootstrap this, DryRun always gets initialized as false and later enabled if needed
	}
	backendCommands.Config = &repoConfig
	prodRunner := git.ProdRunner{
		Config:  repoConfig,
		Backend: backendCommands,
		Frontend: git.FrontendCommands{
			FrontendRunner: newFrontendRunner(newFrontendRunnerArgs{
				omitBranchNames:  args.OmitBranchNames,
				dryRun:           args.DryRun,
				getCurrentBranch: backendCommands.CurrentBranch,
				stats:            &commandsCounter,
			}),
			SetCachedCurrentBranch: backendCommands.CurrentBranchCache.Set,
		},
		CommandsCounter: &commandsCounter,
	}
	if args.DryRun {
		prodRunner.Config.DryRun = true
	}
	rootDir := backendCommands.RootDirectory()
	if args.ValidateGitRepo {
		if rootDir.IsEmpty() {
			err = errors.New(messages.RepoOutside)
			return nil, err
		}
	}
	isOffline, err := repoConfig.IsOffline()
	if err != nil {
		return nil, err
	}
	if args.ValidateIsOnline && isOffline {
		err = errors.New(messages.OfflineNotAllowed)
		return nil, err
	}
	if args.ValidateGitRepo {
		var currentDirectory string
		currentDirectory, err = os.Getwd()
		if err != nil {
			err = errors.New(messages.DirCurrentProblem)
			return nil, err
		}
		if currentDirectory != rootDir.String() {
			err = prodRunner.Frontend.NavigateToDir(rootDir)
		}
	}
	return &OpenRepoResult{
		Runner:         prodRunner,
		RootDir:        rootDir,
		IsOffline:      isOffline,
		ConfigSnapshot: configSnapshot,
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
	Runner         git.ProdRunner
	RootDir        domain.RepoRootDir
	IsOffline      bool
	ConfigSnapshot undo.ConfigSnapshot
}

// newFrontendRunner provides a FrontendRunner instance that behaves according to the given configuration.
func newFrontendRunner(args newFrontendRunnerArgs) git.FrontendRunner {
	if args.dryRun {
		return &subshell.FrontendDryRunner{
			GetCurrentBranch: args.getCurrentBranch,
			OmitBranchNames:  args.omitBranchNames,
			CommandsCounter:  &args.stats,
		}
	}
	return &subshell.FrontendRunner{
		GetCurrentBranch: args.getCurrentBranch,
		OmitBranchNames:  args.omitBranchNames,
		Stats:            args.stats,
	}
}

type newFrontendRunnerArgs struct {
	omitBranchNames  bool
	dryRun           bool
	getCurrentBranch subshell.GetCurrentBranchFunc
	counter          *gohacks.Counter
}
