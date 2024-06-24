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

func OpenRepo(args OpenRepoArgs) (OpenRepoResult, error) {
	commandsCounter := NewMutable(new(gohacks.Counter))
	backendRunner := subshell.BackendRunner{
		Dir:             None[string](),
		CommandsCounter: commandsCounter,
		Verbose:         args.Verbose,
	}
	gitCommands := git.Commands{
		CurrentBranchCache: &cache.LocalBranchWithPrevious{},
		RemotesCache:       &cache.Remotes{},
	}
	gitVersionMajor, gitVersionMinor, err := gitCommands.Version(backendRunner)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	err = validate.HasAcceptableGitVersion(gitVersionMajor, gitVersionMinor)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	configGitAccess := gitconfig.Access{Runner: backendRunner}
	globalSnapshot, globalConfig, err := configGitAccess.LoadGlobal(true)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	localSnapshot, localConfig, err := configGitAccess.LoadLocal(true)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	configSnapshot := undoconfig.ConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	}
	configFile, err := configfile.Load()
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	unvalidatedConfig, finalMessages := config.NewUnvalidatedConfig(config.NewUnvalidatedConfigArgs{
		Access:       configGitAccess,
		ConfigFile:   configFile,
		DryRun:       args.DryRun,
		GlobalConfig: globalConfig,
		LocalConfig:  localConfig,
	})
	frontEndRunner := newFrontendRunner(newFrontendRunnerArgs{
		backend:          backendRunner,
		counter:          commandsCounter,
		dryRun:           args.DryRun,
		getCurrentBranch: gitCommands.CurrentBranch,
		omitBranchNames:  args.OmitBranchNames,
		printCommands:    args.PrintCommands,
	})
	rootDir, hasRootDir := gitCommands.RootDirectory(backendRunner).Get()
	if args.ValidateGitRepo {
		if !hasRootDir {
			err = errors.New(messages.RepoOutside)
			return emptyOpenRepoResult(), err
		}
	}
	isOffline := unvalidatedConfig.Config.Value.Offline
	if args.ValidateIsOnline && isOffline.Bool() {
		err = errors.New(messages.OfflineNotAllowed)
		return emptyOpenRepoResult(), err
	}
	if args.ValidateGitRepo {
		var currentDirectory string
		currentDirectory, err = os.Getwd()
		if err != nil {
			err = errors.New(messages.DirCurrentProblem)
			return emptyOpenRepoResult(), err
		}
		if currentDirectory != rootDir.String() {
			err = gitCommands.NavigateToDir(rootDir)
		}
	}
	return OpenRepoResult{
		Backend:           backendRunner,
		CommandsCounter:   commandsCounter,
		ConfigSnapshot:    configSnapshot,
		FinalMessages:     finalMessages,
		Frontend:          frontEndRunner,
		Git:               gitCommands,
		IsOffline:         isOffline,
		RootDir:           rootDir,
		UnvalidatedConfig: unvalidatedConfig,
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
	Backend           gitdomain.RunnerQuerier
	CommandsCounter   Mutable[gohacks.Counter]
	ConfigSnapshot    undoconfig.ConfigSnapshot
	FinalMessages     stringslice.Collector
	Frontend          gitdomain.Runner
	Git               git.Commands
	IsOffline         configdomain.Offline
	RootDir           gitdomain.RepoRootDir
	UnvalidatedConfig config.UnvalidatedConfig
}

func emptyOpenRepoResult() OpenRepoResult {
	return OpenRepoResult{} //exhaustruct:ignore
}

// newFrontendRunner provides a FrontendRunner instance that behaves according to the given configuration.
func newFrontendRunner(args newFrontendRunnerArgs) gitdomain.Runner { //nolint:ireturn
	if args.dryRun {
		return &subshell.FrontendDryRunner{
			Backend:          args.backend,
			GetCurrentBranch: args.getCurrentBranch,
			OmitBranchNames:  args.omitBranchNames,
			PrintCommands:    args.printCommands,
			CommandsCounter:  args.counter,
		}
	}
	return &subshell.FrontendRunner{
		Backend:          args.backend,
		GetCurrentBranch: args.getCurrentBranch,
		OmitBranchNames:  args.omitBranchNames,
		PrintCommands:    args.printCommands,
		CommandsCounter:  args.counter,
	}
}

type newFrontendRunnerArgs struct {
	backend          gitdomain.Querier
	counter          Mutable[gohacks.Counter]
	dryRun           bool
	getCurrentBranch subshell.GetCurrentBranchFunc
	omitBranchNames  bool
	printCommands    bool
}
