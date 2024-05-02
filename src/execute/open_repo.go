package execute

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/configfile"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/gohacks/cache"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/subshell"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
)

func OpenRepo(args OpenRepoArgs) (*OpenRepoResult, error) {
	commandsCounter := gohacks.Counter{}
	backendRunner := subshell.BackendRunner{
		Dir:             None[string](),
		CommandsCounter: &commandsCounter,
		Verbose:         args.Verbose,
	}
	backendCommands := git.BackendCommands{
		Runner:             backendRunner,
		DryRun:             args.DryRun,
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
	configGitAccess := gitconfig.Access{Runner: backendRunner}
	globalSnapshot, globalConfig, err := configGitAccess.LoadGlobal(true)
	if err != nil {
		return nil, err
	}
	localSnapshot, localConfig, err := configGitAccess.LoadLocal(true)
	if err != nil {
		return nil, err
	}
	configSnapshot := undoconfig.ConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	}
	configFile, err := configfile.Load()
	if err != nil {
		return nil, err
	}
	config, finalMessages := config.NewConfig(config.NewConfigArgs{
		ConfigFile:   configFile,
		DryRun:       args.DryRun,
		GlobalConfig: globalConfig,
		LocalConfig:  localConfig,
		Runner:       backendRunner,
	})
	frontEndRunner := newFrontendRunner(newFrontendRunnerArgs{
		counter:          &commandsCounter,
		dryRun:           args.DryRun,
		getCurrentBranch: backendCommands.CurrentBranch,
		omitBranchNames:  args.OmitBranchNames,
		printCommands:    args.PrintCommands,
	})
	frontEndCommands := git.FrontendCommands{
		Runner:                 frontEndRunner,
		SetCachedCurrentBranch: backendCommands.CurrentBranchCache.Set,
	}
	prodRunner := git.ProdRunner{
		Config:          config,
		Backend:         backendCommands,
		Frontend:        frontEndCommands,
		CommandsCounter: &commandsCounter,
		FinalMessages:   finalMessages,
	}
	rootDir := backendCommands.RootDirectory()
	if args.ValidateGitRepo {
		if rootDir.IsEmpty() {
			err = errors.New(messages.RepoOutside)
			return nil, err
		}
	}
	isOffline := config.Config.Offline
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
		Backend:         backendCommands,
		CommandsCounter: &commandsCounter,
		Config:          config,
		ConfigSnapshot:  configSnapshot,
		FinalMessages:   finalMessages,
		Frontend:        frontEndCommands,
		IsOffline:       isOffline,
		RootDir:         rootDir,
	}, err
}

type OpenRepoArgs struct {
	DryRun           bool
	OmitBranchNames  bool
	PrintCommands    bool
	ValidateGitRepo  bool
	ValidateIsOnline bool
	Verbose          bool
}

type OpenRepoResult struct {
	Backend         git.BackendCommands
	CommandsCounter *gohacks.Counter
	Config          *config.Config
	ConfigSnapshot  undoconfig.ConfigSnapshot
	FinalMessages   *stringslice.Collector
	Frontend        git.FrontendCommands
	IsOffline       configdomain.Offline
	RootDir         gitdomain.RepoRootDir
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
	counter          *gohacks.Counter
	dryRun           bool
	getCurrentBranch subshell.GetCurrentBranchFunc
	omitBranchNames  bool
	printCommands    bool
}
