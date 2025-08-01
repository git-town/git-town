package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/configinterpreter"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	contributeDesc = "Stop syncing some feature branches with their parents"
	contributeHelp = `
Marks the given local branches as contribution.
If no branch is provided, marks the current branch.

Contribution branches are useful when you assist other developers
and make commits to their branch,
but want the other developers to manage the branch
including syncing it with its parent and shipping it.

On a contribution branch, "git town sync"
- pulls down updates from the tracking branch (always via rebase)
- pushes your local commits to the tracking branch
- does not pull updates from the parent branch
`
)

func contributeCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "contribute [branches]",
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDTypes,
		Short:   contributeDesc,
		Long:    cmdhelpers.Long(contributeDesc, contributeHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: None[configdomain.AutoResolve](),
				DryRun:      None[configdomain.DryRun](),
				Verbose:     verbose,
			})
			return executeContribute(args, cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeContribute(args []string, cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, err := determineContributeData(args, repo)
	if err != nil {
		return err
	}
	if err = validateContributeData(data, repo); err != nil {
		return err
	}
	branchNames := data.branchesToMark.Keys()
	if err = gitconfig.SetBranchTypeOverride(repo.Backend, configdomain.BranchTypeContributionBranch, branchNames...); err != nil {
		return err
	}
	printContributeBranches(branchNames)
	if branchToCheckout, hasBranchToCheckout := data.branchToCheckout.Get(); hasBranchToCheckout {
		if err = repo.Git.CheckoutBranch(repo.Frontend, branchToCheckout, false); err != nil {
			return err
		}
	}
	return configinterpreter.Finished(configinterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: Some(data.beginBranchesSnapshot),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "contribute",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       data.branchesToMark.Keys().BranchNames(),
		Verbose:               repo.UnvalidatedConfig.NormalConfig.Verbose,
	})
}

type contributeData struct {
	beginBranchesSnapshot gitdomain.BranchesSnapshot
	branchInfos           gitdomain.BranchInfos
	branchToCheckout      Option[gitdomain.LocalBranchName]
	branchesToMark        configdomain.BranchesAndTypes
}

func printContributeBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.BranchIsNowContribution, branch)
	}
}

func determineContributeData(args []string, repo execute.OpenRepoResult) (contributeData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return contributeData{}, err
	}
	branchesToMakeContribution, branchToCheckout, err := config.BranchesToMark(args, branchesSnapshot, repo.UnvalidatedConfig)
	return contributeData{
		beginBranchesSnapshot: branchesSnapshot,
		branchInfos:           branchesSnapshot.Branches,
		branchToCheckout:      branchToCheckout,
		branchesToMark:        branchesToMakeContribution,
	}, err
}

func validateContributeData(data contributeData, repo execute.OpenRepoResult) error {
	for branchName, branchType := range data.branchesToMark {
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotMakeContribution)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotMakeContribution)
		case configdomain.BranchTypeContributionBranch:
			return fmt.Errorf(messages.BranchIsAlreadyContribution, branchName)
		case
			configdomain.BranchTypeFeatureBranch,
			configdomain.BranchTypeObservedBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypePrototypeBranch:
		}
		hasLocalBranch := data.branchInfos.HasLocalBranch(branchName)
		hasRemoteBranch := data.branchInfos.HasMatchingTrackingBranchFor(branchName, repo.UnvalidatedConfig.NormalConfig.DevRemote)
		if !hasLocalBranch && !hasRemoteBranch {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
		if hasLocalBranch && !hasRemoteBranch {
			return fmt.Errorf(messages.ContributeBranchIsLocal, branchName)
		}
	}
	return nil
}
