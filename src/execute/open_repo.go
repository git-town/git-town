package execute

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/git"
	"github.com/git-town/git-town/v10/src/gohacks"
	"github.com/git-town/git-town/v10/src/gohacks/cache"
	"github.com/git-town/git-town/v10/src/gohacks/stringslice"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/subshell"
	"github.com/git-town/git-town/v10/src/undo"
	"github.com/git-town/git-town/v10/src/validate"
)

func OpenRepo(args OpenRepoArgs) (*OpenRepoResult, error) {
	commandsCounter := gohacks.Counter{}
	backendRunner := subshell.BackendRunner{
		Dir:             nil,
		CommandsCounter: &commandsCounter,
		Verbose:         args.Verbose,
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
				printCommands:    args.PrintCommands,
				dryRun:           args.DryRun,
				getCurrentBranch: backendCommands.CurrentBranch,
				counter:          &commandsCounter,
			}),
			SetCachedCurrentBranch: backendCommands.CurrentBranchCache.Set,
		},
		CommandsCounter: &commandsCounter,
		FinalMessages:   &stringslice.Collector{},
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
	Verbose          bool
	DryRun           bool
	OmitBranchNames  bool
	PrintCommands    bool
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
			PrintCommands:    args.printCommands,
			CommandsCounter:  args.counter,
		}
	}
	return &subshell.FrontendRunner{
		GetCurrentBranch: args.getCurrentBranch,
		OmitBranchNames:  args.omitBranchNames,
		PrintCommands:    args.printCommands,
		CommandsCounter:  args.counter,
	}
}

type newFrontendRunnerArgs struct {
	omitBranchNames  bool
	printCommands    bool
	dryRun           bool
	getCurrentBranch subshell.GetCurrentBranchFunc
	counter          *gohacks.Counter
}
