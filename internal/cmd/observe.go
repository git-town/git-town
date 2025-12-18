package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/mapstools"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/configinterpreter"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	observeDesc = "Stop your contributions to some feature branches"
	observeHelp = `
Marks the given local branches as observed.
If no branch is provided, observes the current branch.

Observed branches are useful when you assist other developers
and make local changes to try out ideas,
but want the other developers to implement and commit all official changes.

On an observed branch, "git town sync"
- pulls down updates from the tracking branch (always via rebase)
- does not push your local commits to the tracking branch
- does not pull updates from the parent branch
`
)

func observeCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "observe [branches]",
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDTypes,
		Short:   observeDesc,
		Long:    cmdhelpers.Long(observeDesc, observeHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeObserve(args, cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeObserve(args []string, cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, err := determineObserveData(args, repo)
	if err != nil {
		return err
	}
	err = validateObserveData(data, repo)
	if err != nil {
		return err
	}
	branchNames := data.branchesToObserve.Keys()
	if err = gitconfig.SetBranchTypeOverride(repo.Backend, configdomain.BranchTypeObservedBranch, branchNames...); err != nil {
		return err
	}
	printObservedBranches(branchNames)
	if checkout, hasCheckout := data.checkout.Get(); hasCheckout {
		if err = repo.Git.CheckoutBranch(repo.Frontend, checkout, false); err != nil {
			return err
		}
	}
	return configinterpreter.Finished(configinterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: Some(data.branchesSnapshot),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "observe",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       branchNames.BranchNames(),
		Verbose:               repo.UnvalidatedConfig.NormalConfig.Verbose,
	})
}

type observeData struct {
	branchInfos       gitdomain.BranchInfos
	branchesSnapshot  gitdomain.BranchesSnapshot
	branchesToObserve configdomain.BranchesAndTypes
	checkout          Option[gitdomain.LocalBranchName]
}

func printObservedBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.BranchIsNowObserved, branch)
	}
}

func determineObserveData(args []string, repo execute.OpenRepoResult) (observeData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return observeData{}, err
	}
	if branchesSnapshot.DetachedHead {
		return observeData{}, errors.New(messages.ObserveDetachedHead)
	}
	branchesToObserve, branchToCheckout, err := config.BranchesToMark(args, branchesSnapshot, repo.UnvalidatedConfig)
	return observeData{
		branchInfos:       branchesSnapshot.Branches,
		branchesSnapshot:  branchesSnapshot,
		branchesToObserve: branchesToObserve,
		checkout:          branchToCheckout,
	}, err
}

func validateObserveData(data observeData, repo execute.OpenRepoResult) error {
	for branchName, branchType := range mapstools.SortedKeyValues(data.branchesToObserve) {
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotObserve)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotObserve)
		case configdomain.BranchTypeObservedBranch:
			repo.FinalMessages.Add(fmt.Sprintf(messages.BranchIsAlreadyObserved, branchName))
		case
			configdomain.BranchTypeFeatureBranch,
			configdomain.BranchTypeContributionBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypePrototypeBranch:
		}
		hasLocalBranch := data.branchesSnapshot.Branches.HasLocalBranch(branchName)
		hasRemoteBranch := data.branchesSnapshot.Branches.HasMatchingTrackingBranchFor(branchName, repo.UnvalidatedConfig.NormalConfig.DevRemote)
		if !hasLocalBranch && !hasRemoteBranch {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
		if hasLocalBranch && !hasRemoteBranch {
			return fmt.Errorf(messages.ObserveBranchIsLocal, branchName)
		}
	}
	return nil
}
