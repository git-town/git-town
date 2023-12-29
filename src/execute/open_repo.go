package execute

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v11/src/config"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/gohacks/cache"
	"github.com/git-town/git-town/v11/src/gohacks/stringslice"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/subshell"
	"github.com/git-town/git-town/v11/src/undo/undoconfig"
	"github.com/git-town/git-town/v11/src/validate"
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
		DryRun:             args.DryRun,
		Config:             nil, // initializing to nil here to validate the Git version before running any Git commands, setting to the correct value after that is done
		CurrentBranchCache: &cache.LocalBranchWithPrevious{},
		RemotesCache:       &cache.Remotes{},
	}
	gitVersionMajor, gitVersionMinor, err := backendCommands.Version()
	if err != nil {
		return nil, err
	}
	err = validate.HasAcceptableGitVersion(gitVersionMajor, gitVersionMinor)
	if err != nil {
		return nil, err
	}
	configGitAccess := configdomain.Access{Runner: backendRunner}
	globalCache, globalConfig, err := configGitAccess.LoadCache(true)
	if err != nil {
		return nil, err
	}
	localCache, localConfig, err := configGitAccess.LoadCache(false)
	if err != nil {
		return nil, err
	}
	configSnapshot := undoconfig.ConfigSnapshot{
		Global: globalCache,
		Local:  localCache,
	}
	gitTown, err := config.NewGitTown(globalConfig, localConfig, args.DryRun, backendRunner)
	if err != nil {
		return nil, err
	}
	backendCommands.Config = gitTown
	prodRunner := git.ProdRunner{
		Config:  gitTown,
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
	rootDir := backendCommands.RootDirectory()
	if args.ValidateGitRepo {
		if rootDir.IsEmpty() {
			err = errors.New(messages.RepoOutside)
			return nil, err
		}
	}
	isOffline := gitTown.FullConfig.Offline
	if args.ValidateIsOnline && isOffline.Bool() {
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
		Runner:         &prodRunner,
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
	Runner         *git.ProdRunner
	RootDir        gitdomain.RepoRootDir
	IsOffline      configdomain.Offline
	ConfigSnapshot undoconfig.ConfigSnapshot
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
