package execute

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configfile"
	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/gohacks/cache"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func OpenRepo(args OpenRepoArgs) (OpenRepoResult, error) {
	commandsCounter := NewMutable(new(gohacks.Counter))
	if args.Verbose {
		fmt.Println("Git Town " + config.GitTownVersion)
		fmt.Println("OS:", runtime.GOOS)
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", "ver")
		} else {
			cmd = exec.Command("uname", "-a")
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}
	backendRunner := subshell.BackendRunner{
		Dir:             None[string](),
		CommandsCounter: commandsCounter,
		Verbose:         args.Verbose,
	}
	gitCommands := git.Commands{
		CurrentBranchCache: &cache.WithPrevious[gitdomain.LocalBranchName]{},
		RemotesCache:       &cache.Cache[gitdomain.Remotes]{},
	}
	gitVersion, err := gitCommands.GitVersion(backendRunner)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	if !gitVersion.IsMinimumRequiredGitVersion() {
		return emptyOpenRepoResult(), errors.New(messages.GitVersionTooLow)
	}
	rootDir, hasRootDir := gitCommands.RootDirectory(backendRunner).Get()
	if args.ValidateGitRepo {
		if !hasRootDir {
			return emptyOpenRepoResult(), errors.New(messages.RepoOutside)
		}
	}
	configGitAccess := gitconfig.Access{Shell: backendRunner}
	globalSnapshot, err := configGitAccess.Load(Some(configdomain.ConfigScopeGlobal), configdomain.UpdateOutdatedYes)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	localSnapshot, err := configGitAccess.Load(Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedYes)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	unscopedSnapshot, err := configGitAccess.Load(None[configdomain.ConfigScope](), configdomain.UpdateOutdatedNo)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	unscopedConfig, err := configdomain.NewPartialConfigFromSnapshot(unscopedSnapshot, true, configGitAccess.RemoveLocalConfigValue)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	configSnapshot := undoconfig.ConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	}
	finalMessages := stringslice.NewCollector()
	configFile, err := configfile.Load(rootDir, configfile.FileName, finalMessages)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	if configFile.IsNone() {
		configFile, err = configfile.Load(rootDir, configfile.AlternativeFileName, finalMessages)
		if err != nil {
			return emptyOpenRepoResult(), err
		}
	}
	unvalidatedConfig := config.NewUnvalidatedConfig(config.NewUnvalidatedConfigArgs{
		Access:        configGitAccess,
		ConfigFile:    configFile,
		DryRun:        args.DryRun,
		EnvConfig:     envconfig.Load(),
		FinalMessages: finalMessages,
		GitConfig:     unscopedConfig,
		GitVersion:    gitVersion,
	})
	frontEndRunner := newFrontendRunner(newFrontendRunnerArgs{
		backend:          backendRunner,
		counter:          commandsCounter,
		dryRun:           args.DryRun,
		getCurrentBranch: gitCommands.CurrentBranch,
		printBranchNames: args.PrintBranchNames,
		printCommands:    args.PrintCommands,
	})
	isOffline := unvalidatedConfig.NormalConfig.Offline
	if args.ValidateIsOnline && isOffline.IsOffline() {
		return emptyOpenRepoResult(), errors.New(messages.OfflineNotAllowed)
	}
	if args.ValidateGitRepo {
		var currentDirectory string
		currentDirectory, err = os.Getwd()
		if err != nil {
			return emptyOpenRepoResult(), errors.New(messages.DirCurrentProblem)
		}
		if currentDirectory != rootDir.String() {
			err = gitCommands.ChangeDir(rootDir)
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
	DryRun           configdomain.DryRun
	PrintBranchNames bool
	PrintCommands    bool
	ValidateGitRepo  bool
	ValidateIsOnline bool
	Verbose          configdomain.Verbose
}

type OpenRepoResult struct {
	Backend           subshelldomain.RunnerQuerier
	CommandsCounter   Mutable[gohacks.Counter]
	ConfigSnapshot    undoconfig.ConfigSnapshot
	FinalMessages     stringslice.Collector
	Frontend          subshelldomain.Runner
	Git               git.Commands
	IsOffline         configdomain.Offline
	RootDir           gitdomain.RepoRootDir
	UnvalidatedConfig config.UnvalidatedConfig
}

func emptyOpenRepoResult() OpenRepoResult {
	return OpenRepoResult{} //exhaustruct:ignore
}

// newFrontendRunner provides a FrontendRunner instance that behaves according to the given configuration.
func newFrontendRunner(args newFrontendRunnerArgs) subshelldomain.Runner { //nolint:ireturn
	if args.dryRun {
		return &subshell.FrontendDryRunner{
			Backend:          args.backend,
			GetCurrentBranch: args.getCurrentBranch,
			PrintBranchNames: args.printBranchNames,
			PrintCommands:    args.printCommands,
			CommandsCounter:  args.counter,
		}
	}
	return &subshell.FrontendRunner{
		Backend:          args.backend,
		GetCurrentBranch: args.getCurrentBranch,
		PrintBranchNames: args.printBranchNames,
		PrintCommands:    args.printCommands,
		CommandsCounter:  args.counter,
	}
}

type newFrontendRunnerArgs struct {
	backend          subshelldomain.Querier
	counter          Mutable[gohacks.Counter]
	dryRun           configdomain.DryRun
	getCurrentBranch subshell.GetCurrentBranchFunc
	printBranchNames bool
	printCommands    bool
}
