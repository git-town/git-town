package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v13/src/cli/flags"
	"github.com/git-town/git-town/v13/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v13/src/config"
	"github.com/git-town/git-town/v13/src/config/commandconfig"
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/execute"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/messages"
	"github.com/git-town/git-town/v13/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v13/src/vm/interpreter/config"
	"github.com/spf13/cobra"
)

const observeDesc = "Stops your contributions to some feature branches"

const observeHelp = `
Marks the given local branches as observed.
If no branch is provided, observes the current branch.

Observed branches are useful when you assist other developers
and make local changes to try out ideas,
but want the other developers to implement and commit all official changes.

On an observed branch, "git sync"
- pulls down updates from the tracking branch (always via rebase)
- does not push your local commits to the tracking branch
- does not pull updates from the parent branch
`

func observeCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "observe [branches]",
		Args:    cobra.ArbitraryArgs,
		GroupID: "types",
		Short:   observeDesc,
		Long:    cmdhelpers.Long(observeDesc, observeHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeObserve(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeObserve(args []string, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	config, err := determineObserveConfig(args, repo)
	if err != nil {
		return err
	}
	err = validateObserveConfig(config)
	if err != nil {
		return err
	}
	branchNames := config.branchesToObserve.Keys()
	if err = repo.Runner.Config.AddToObservedBranches(branchNames...); err != nil {
		return err
	}
	if err = removeNonObserveBranchTypes(config.branchesToObserve, repo.Runner.Config); err != nil {
		return err
	}
	printObservedBranches(branchNames)
	if !config.checkout.IsEmpty() {
		if err = repo.Runner.Frontend.CheckoutBranch(config.checkout); err != nil {
			return err
		}
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		BeginConfigSnapshot: repo.ConfigSnapshot,
		Command:             "observe",
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		RootDir:             repo.RootDir,
		Runner:              repo.Runner,
		Verbose:             verbose,
	})
}

type observeConfig struct {
	allBranches       gitdomain.BranchInfos
	branchesToObserve commandconfig.BranchesAndTypes
	checkout          gitdomain.LocalBranchName
}

func printObservedBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.ObservedBranchIsNowObserved, branch)
	}
}

func removeNonObserveBranchTypes(branches map[gitdomain.LocalBranchName]configdomain.BranchType, config *config.Config) error {
	for branchName, branchType := range branches {
		switch branchType {
		case configdomain.BranchTypeContributionBranch:
			if err := config.RemoveFromContributionBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeParkedBranch:
			if err := config.RemoveFromParkedBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		}
	}
	return nil
}

func determineObserveConfig(args []string, repo *execute.OpenRepoResult) (observeConfig, error) {
	branchesSnapshot, err := repo.Runner.Backend.BranchesSnapshot()
	if err != nil {
		return observeConfig{}, err
	}
	branchesToObserve := commandconfig.BranchesAndTypes{}
	checkout := gitdomain.EmptyLocalBranchName()
	switch len(args) {
	case 0:
		branchesToObserve.Add(branchesSnapshot.Active, &repo.Runner.Config.FullConfig)
	case 1:
		branch := gitdomain.NewLocalBranchName(args[0])
		branchesToObserve.Add(branch, &repo.Runner.Config.FullConfig)
		branchInfo := branchesSnapshot.Branches.FindByRemoteName(branch.TrackingBranch())
		if branchInfo.SyncStatus == gitdomain.SyncStatusRemoteOnly {
			checkout = branch
		}
	default:
		branchesToObserve.AddMany(gitdomain.NewLocalBranchNames(args...), &repo.Runner.Config.FullConfig)
	}
	return observeConfig{
		allBranches:       branchesSnapshot.Branches,
		branchesToObserve: branchesToObserve,
		checkout:          checkout,
	}, nil
}

func validateObserveConfig(config observeConfig) error {
	for branchName, branchType := range config.branchesToObserve {
		if !config.allBranches.HasLocalBranch(branchName) && !config.allBranches.HasMatchingTrackingBranchFor(branchName) {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotObserve)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotObserve)
		case configdomain.BranchTypeObservedBranch:
			return fmt.Errorf(messages.BranchIsAlreadyObserved, branchName)
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeContributionBranch, configdomain.BranchTypeParkedBranch:
		}
	}
	return nil
}
