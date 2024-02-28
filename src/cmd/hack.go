package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config"
	"github.com/git-town/git-town/v12/src/config/commandconfig"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/config"
	fullInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/full"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/spf13/cobra"
)

const hackDesc = "Creates a new feature branch off the main development branch"

const hackHelp = `
Syncs the main branch, forks a new feature branch with the given name off the main branch, pushes the new feature branch to origin (if and only if "push-new-branches" is true), and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`

func hackCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     "hack <branch>",
		GroupID: "basic",
		Args:    cobra.ArbitraryArgs,
		Short:   hackDesc,
		Long:    cmdhelpers.Long(hackDesc, hackHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeHack(args, readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeHack(args []string, dryRun, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineHackConfig(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	if config.appendConfig != nil {
		return createBranch(repo, config.appendConfig, initialBranchesSnapshot, initialStashSize, dryRun, verbose)
	}
	if config.makeFeatureConfig != nil {
		return makeFeatureBranch(makeFeatureBranchArgs{
			config:              repo.Runner.Config,
			beginConfigSnapshot: repo.ConfigSnapshot,
			makeFeatureConfig:   config.makeFeatureConfig,
			rootDir:             repo.RootDir,
			runner:              repo.Runner,
			verbose:             verbose,
		})
	}
	panic("both config arms were nil")
}

type hackConfig struct {
	appendConfig      *appendConfig
	makeFeatureConfig *makeFeatureConfig
}

// this configuration is for when "git hack" is used to make contribution, observed, or parked branches feature branches
type makeFeatureConfig struct {
	targetBranches commandconfig.BranchesAndTypes
}

func createBranch(repo *execute.OpenRepoResult, config *appendConfig, initialBranchesSnapshot gitdomain.BranchesSnapshot, initialStashSize gitdomain.StashSize, dryRun, verbose bool) error {
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               "hack",
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            appendProgram(config),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               nil,
		DialogTestInputs:        &config.dialogTestInputs,
		FullConfig:              config.FullConfig,
		HasOpenChanges:          config.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     repo.Runner,
		RunState:                &runState,
		Verbose:                 verbose,
	})
}

func determineHackConfig(args []string, repo *execute.OpenRepoResult, dryRun, verbose bool) (*hackConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	fc := execute.FailureCollector{}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	branchesSnapshot, stashSize, repoStatus, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FullConfig:            &repo.Runner.Config.FullConfig,
		HandleUnfinishedState: true,
		Repo:                  repo,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	targetBranches := gitdomain.NewLocalBranchNames(args...)
	if len(targetBranches) == 0 {
		return &hackConfig{
			appendConfig: nil,
			makeFeatureConfig: &makeFeatureConfig{
				targetBranches: commandconfig.NewBranchesAndTypes(gitdomain.LocalBranchNames{branchesSnapshot.Active}, repo.Runner.Config.FullConfig),
			},
		}, branchesSnapshot, stashSize, false, nil
	}
	if len(targetBranches) > 0 && branchesSnapshot.Branches.HasLocalBranches(targetBranches) {
		return &hackConfig{
			appendConfig: nil,
			makeFeatureConfig: &makeFeatureConfig{
				targetBranches: commandconfig.NewBranchesAndTypes(targetBranches, repo.Runner.Config.FullConfig),
			},
		}, branchesSnapshot, stashSize, false, nil
	}
	if len(targetBranches) > 1 {
		return nil, branchesSnapshot, stashSize, false, errors.New(messages.HackTooManyArguments)
	}
	targetBranch := targetBranches[0]
	remotes := fc.Remotes(repo.Runner.Backend.Remotes())
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	branchNamesToSync := gitdomain.LocalBranchNames{repo.Runner.Config.FullConfig.MainBranch}
	branchesToSync := fc.BranchInfos(branchesSnapshot.Branches.Select(branchNamesToSync))
	return &hackConfig{
		appendConfig: &appendConfig{
			FullConfig:                &repo.Runner.Config.FullConfig,
			allBranches:               branchesSnapshot.Branches,
			branchesToSync:            branchesToSync,
			dialogTestInputs:          dialogTestInputs,
			dryRun:                    dryRun,
			hasOpenChanges:            repoStatus.OpenChanges,
			initialBranch:             branchesSnapshot.Active,
			newBranchParentCandidates: gitdomain.LocalBranchNames{repo.Runner.Config.FullConfig.MainBranch},
			parentBranch:              repo.Runner.Config.FullConfig.MainBranch,
			previousBranch:            previousBranch,
			remotes:                   remotes,
			targetBranch:              targetBranch,
		},
		makeFeatureConfig: nil,
	}, branchesSnapshot, stashSize, false, fc.Err
}

func makeFeatureBranch(args makeFeatureBranchArgs) error {
	err := validateMakeFeatureConfig(args.makeFeatureConfig)
	if err != nil {
		return err
	}
	for branchName, branchType := range args.makeFeatureConfig.targetBranches {
		switch branchType {
		case configdomain.BranchTypeContributionBranch:
			err = args.config.RemoveFromContributionBranches(branchName)
		case configdomain.BranchTypeObservedBranch:
			err = args.config.RemoveFromObservedBranches(branchName)
		case configdomain.BranchTypeParkedBranch:
			err = args.config.RemoveFromParkedBranches(branchName)
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
			panic(fmt.Sprintf("unchecked branch type: %s", branchType))
		}
		if err != nil {
			return err
		}
		fmt.Printf(messages.HackBranchIsNowFeature, branchName)
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		BeginConfigSnapshot: args.beginConfigSnapshot,
		Command:             "observe",
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		RootDir:             args.rootDir,
		Runner:              args.runner,
		Verbose:             args.verbose,
	})
}

type makeFeatureBranchArgs struct {
	config              *config.Config
	beginConfigSnapshot undoconfig.ConfigSnapshot
	makeFeatureConfig   *makeFeatureConfig
	rootDir             gitdomain.RepoRootDir
	runner              *git.ProdRunner
	verbose             bool
}

func validateMakeFeatureConfig(config *makeFeatureConfig) error {
	for branchName, branchType := range config.targetBranches {
		switch branchType {
		case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch:
			return nil
		case configdomain.BranchTypeFeatureBranch:
			return fmt.Errorf(messages.HackBranchIsAlreadyFeature, branchName)
		case configdomain.BranchTypeMainBranch:
			return fmt.Errorf(messages.HackCannotFeatureMainBranch, branchName)
		case configdomain.BranchTypePerennialBranch:
			return fmt.Errorf(messages.HackCannotFeaturePerennialBranch, branchName)
		}
		panic(fmt.Sprintf("unhandled branch type: %s", branchType))
	}
	return nil
}
