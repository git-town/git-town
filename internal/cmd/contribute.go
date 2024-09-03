package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	configInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/config"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const contributeDesc = "Stop syncing some feature branches with their parents"

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

func executeContribute(args []string, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, err := determineContributeData(args, repo)
	if err != nil {
		return err
	}
	if err = validateContributeData(data); err != nil {
		return err
	}
	branchNames := data.branchesToMark.Keys()
	if err = repo.UnvalidatedConfig.AddToContributionBranches(branchNames...); err != nil {
		return err
	}
	if err = removeNonContributionBranchTypes(data.branchesToMark, repo.UnvalidatedConfig); err != nil {
		return err
	}
	printContributeBranches(branchNames)
	if branchToCheckout, hasBranchToCheckout := data.branchToCheckout.Get(); hasBranchToCheckout {
		if err = repo.Git.CheckoutBranch(repo.Frontend, branchToCheckout, false); err != nil {
			return err
		}
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: Some(data.beginBranchesSnapshot),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "contribute",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       data.branchesToMark.Keys().BranchNames(),
		Verbose:               verbose,
	})
}

type contributeData struct {
	allBranches           gitdomain.BranchInfos
	beginBranchesSnapshot gitdomain.BranchesSnapshot
	branchToCheckout      Option[gitdomain.LocalBranchName]
	branchesToMark        configdomain.BranchesAndTypes
}

func printContributeBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.ContributeBranchIsNowContribution, branch)
	}
}

func removeNonContributionBranchTypes(branches configdomain.BranchesAndTypes, config config.UnvalidatedConfig) error {
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
		case configdomain.BranchTypePrototypeBranch:
			if err := config.RemoveFromPrototypeBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeContributionBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		}
	}
	return nil
}

func determineContributeData(args []string, repo execute.OpenRepoResult) (contributeData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return contributeData{}, err
	}
	branchesToMakeContribution, branchToCheckout, err := execute.BranchesToMark(args, branchesSnapshot, repo.UnvalidatedConfig.Config.Get())
	return contributeData{
		allBranches:           branchesSnapshot.Branches,
		beginBranchesSnapshot: branchesSnapshot,
		branchToCheckout:      branchToCheckout,
		branchesToMark:        branchesToMakeContribution,
	}, err
}

func validateContributeData(data contributeData) error {
	for branchName, branchType := range data.branchesToMark {
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotMakeContribution)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotMakeContribution)
		case configdomain.BranchTypeContributionBranch:
			return fmt.Errorf(messages.BranchIsAlreadyContribution, branchName)
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
		}
		hasLocalBranch := data.allBranches.HasLocalBranch(branchName)
		hasRemoteBranch := data.allBranches.HasMatchingTrackingBranchFor(branchName)
		if !hasLocalBranch && !hasRemoteBranch {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
		if hasLocalBranch && !hasRemoteBranch {
			return fmt.Errorf(messages.ContributeBranchIsLocal, branchName)
		}
	}
	return nil
}
