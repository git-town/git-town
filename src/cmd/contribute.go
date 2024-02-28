package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config"
	"github.com/git-town/git-town/v12/src/config/commandconfig"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/config"
	"github.com/spf13/cobra"
)

const contributeDesc = "Stops syncing some feature branches with their parents"

const contributeHelp = `
Marks the given local branches as contribution.
If no branch is provided, marks the current branch.

Contribution branches are useful when you assist other developers
and make commits to their branch,
but want the other developers to manage the branch
including syncing it with its parent and shipping it.

On a contribution branch, "git sync"
- pulls down updates from the tracking branch (always via rebase)
- pushes your local commits to the tracking branch
- does not pull updates from the parent branch
`

func contributeCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "contribute [branches]",
		Args:    cobra.ArbitraryArgs,
		GroupID: "types",
		Short:   contributeDesc,
		Long:    cmdhelpers.Long(contributeDesc, contributeHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeContribute(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeContribute(args []string, verbose bool) error {
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
	config, err := determineContributeConfig(args, repo)
	if err != nil {
		return err
	}
	err = validateContributeConfig(config)
	if err != nil {
		return err
	}
	contributionKeys := config.branchesToMark.Keys()
	if err = repo.Runner.Config.AddToContributionBranches(contributionKeys...); err != nil {
		return err
	}
	if err = removeNonContributionBranchTypes(config.branchesToMark, repo.Runner.Config); err != nil {
		return err
	}
	printContributeBranches(contributionKeys)
	if !config.checkout.IsEmpty() {
		if err = repo.Runner.Frontend.CheckoutBranch(config.checkout); err != nil {
			return err
		}
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		BeginConfigSnapshot: repo.ConfigSnapshot,
		Command:             "contribute",
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		RootDir:             repo.RootDir,
		Runner:              repo.Runner,
		Verbose:             verbose,
	})
}

type contributeConfig struct {
	allBranches    gitdomain.BranchInfos
	branchesToMark commandconfig.BranchesAndTypes
	checkout       gitdomain.LocalBranchName
}

func printContributeBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.ContributeBranchIsNowContribution, branch)
	}
}

func removeNonContributionBranchTypes(branches commandconfig.BranchesAndTypes, config *config.Config) error {
	for branchName, branchType := range branches {
		switch branchType {
		case configdomain.BranchTypeObservedBranch:
			if err := config.RemoveFromObservedBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeParkedBranch:
			if err := config.RemoveFromParkedBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeContributionBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		}
	}
	return nil
}

func determineContributeConfig(args []string, repo *execute.OpenRepoResult) (contributeConfig, error) {
	branchesSnapshot, err := repo.Runner.Backend.BranchesSnapshot()
	if err != nil {
		return contributeConfig{}, err
	}
	branchesToMark := commandconfig.BranchesAndTypes{}
	checkout := gitdomain.EmptyLocalBranchName()
	switch len(args) {
	case 0:
		branchesToMark.Add(branchesSnapshot.Active, &repo.Runner.Config.FullConfig)
	case 1:
		branch := gitdomain.NewLocalBranchName(args[0])
		branchesToMark.Add(branch, &repo.Runner.Config.FullConfig)
		branchInfo := branchesSnapshot.Branches.FindByRemoteName(branch.TrackingBranch())
		if branchInfo.SyncStatus == gitdomain.SyncStatusRemoteOnly {
			checkout = branch
		}
	default:
		branchesToMark.AddMany(gitdomain.NewLocalBranchNames(args...), &repo.Runner.Config.FullConfig)
	}
	return contributeConfig{
		allBranches:    branchesSnapshot.Branches,
		branchesToMark: branchesToMark,
		checkout:       checkout,
	}, nil
}

func validateContributeConfig(config contributeConfig) error {
	for branchName, branchType := range config.branchesToMark {
		if !config.allBranches.HasLocalBranch(branchName) && !config.allBranches.HasMatchingTrackingBranchFor(branchName) {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotMakeContribution)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotMakeContribution)
		case configdomain.BranchTypeContributionBranch:
			return fmt.Errorf(messages.BranchIsAlreadyContribution, branchName)
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch:
		}
	}
	return nil
}
