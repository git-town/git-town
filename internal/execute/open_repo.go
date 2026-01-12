package execute

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/configfile"
	"github.com/git-town/git-town/v22/internal/config/envconfig"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/cache"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func OpenRepo(args OpenRepoArgs) (OpenRepoResult, error) {
	defaultConfig := config.DefaultNormalConfig()
	envConfig, err := envconfig.Load(envconfig.NewEnvVars(os.Environ()))
	if err != nil {
		return emptyOpenRepoResult(), fmt.Errorf("error loading configuration from environment variables: %w", err)
	}
	commandsCounter := NewMutable(new(gohacks.Counter))
	backendRunner := subshell.BackendRunner{
		Dir:             None[string](),
		CommandsCounter: commandsCounter,
		Verbose:         args.CliConfig.Verbose.Or(envConfig.Verbose).GetOr(defaultConfig.Verbose),
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
	globalSnapshot, err := gitconfig.LoadSnapshot(backendRunner, Some(configdomain.ConfigScopeGlobal), configdomain.UpdateOutdatedYes)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	globalConfig, err := config.NewPartialConfigFromSnapshot(globalSnapshot, true, args.IgnoreUnknown, backendRunner)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	localSnapshot, err := gitconfig.LoadSnapshot(backendRunner, Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedYes)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	localConfig, err := config.NewPartialConfigFromSnapshot(localSnapshot, true, args.IgnoreUnknown, backendRunner)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	unscopedSnapshot, err := gitconfig.LoadSnapshot(backendRunner, None[configdomain.ConfigScope](), configdomain.UpdateOutdatedNo)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	unscopedConfig, err := config.NewPartialConfigFromSnapshot(unscopedSnapshot, true, args.IgnoreUnknown, backendRunner)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	configSnapshot := configdomain.BeginConfigSnapshot{
		Global:   globalSnapshot,
		Local:    localSnapshot,
		Unscoped: unscopedSnapshot,
	}
	finalMessages := stringslice.NewCollector()
	configFile, hasConfigFile, err := configfile.Load(rootDir, configfile.FileName, finalMessages)
	if err != nil {
		return emptyOpenRepoResult(), err
	}
	if !hasConfigFile {
		configFile, hasConfigFile, err = configfile.Load(rootDir, configfile.HiddenFileName, finalMessages)
		if err != nil {
			return emptyOpenRepoResult(), err
		}
	}
	if !hasConfigFile {
		configFile, _, err = configfile.Load(rootDir, configfile.AlternativeFileName, finalMessages)
		if err != nil {
			return emptyOpenRepoResult(), err
		}
	}
	unvalidatedConfig := config.NewUnvalidatedConfig(config.NewUnvalidatedConfigArgs{
		CliConfig:     args.CliConfig,
		ConfigFile:    configFile,
		Defaults:      defaultConfig,
		EnvConfig:     envConfig,
		FinalMessages: finalMessages,
		GitGlobal:     globalConfig,
		GitLocal:      localConfig,
		GitUnscoped:   unscopedConfig,
	})
	backendRunner.Verbose = unvalidatedConfig.NormalConfig.Verbose
	frontEndRunner := newFrontendRunner(newFrontendRunnerArgs{
		backend:          backendRunner,
		counter:          commandsCounter,
		dryRun:           unvalidatedConfig.NormalConfig.DryRun,
		getCurrentBranch: gitCommands.CurrentBranch,
		getCurrentSHA:    gitCommands.CurrentSHA,
		printBranchNames: args.PrintBranchNames,
		printCommands:    args.PrintCommands,
	})
	if unvalidatedConfig.NormalConfig.Verbose {
		fmt.Println("Git Town " + config.GitTownVersion)
		fmt.Println("OS:", runtime.GOOS)
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.CommandContext(context.Background(), "cmd", "/c", "ver")
		} else {
			cmd = exec.CommandContext(context.Background(), "uname", "-a")
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}
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
	CliConfig        configdomain.PartialConfig
	IgnoreUnknown    bool // whether to ignore unknown configuration values
	PrintBranchNames bool // whether Git Town output should contain branch names
	PrintCommands    bool // whether Git Town output should list Git commands that Git Town executes
	ValidateGitRepo  bool // whether to ensure whether the current directory is a Git repository
	ValidateIsOnline bool // whether this Git Town command requires an online connection
}

type OpenRepoResult struct {
	Backend           subshelldomain.RunnerQuerier
	CommandsCounter   Mutable[gohacks.Counter]
	ConfigDir         configdomain.ConfigDirRepo
	ConfigSnapshot    configdomain.BeginConfigSnapshot
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
		GetCurrentSHA:    args.getCurrentSHA,
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
	getCurrentSHA    subshell.GetCurrentSHAFunc
	printBranchNames bool
	printCommands    bool
}
