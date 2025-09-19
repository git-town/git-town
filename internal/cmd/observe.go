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
				AutoResolve:  None[configdomain.AutoResolve](),
				AutoSync:     None[configdomain.AutoSync](),
				Detached:     None[configdomain.Detached](),
				DryRun:       None[configdomain.DryRun](),
				PushBranches: None[configdomain.PushBranches](),
				Stash:        None[configdomain.Stash](),
				Verbose:      verbose,
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
	if err = validateObserveData(data, repo); err != nil {
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
	observeCandidates, branchToCheckout, err := config.BranchesToMark(args, branchesSnapshot, repo.UnvalidatedConfig)
	var branchesToObserve configdomain.BranchesAndTypes
	for branchName, branchType := range observeCandidates {
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return observeData{}, errors.New(messages.MainBranchCannotObserve)
		case configdomain.BranchTypePerennialBranch:
			return observeData{}, errors.New(messages.PerennialBranchCannotObserve)
		case configdomain.BranchTypeObservedBranch:
			repo.FinalMessages.Add(fmt.Sprintf(messages.BranchIsAlreadyObserved, branchName))
		case
			configdomain.BranchTypeFeatureBranch,
			configdomain.BranchTypeContributionBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypePrototypeBranch:
			hasLocalBranch := branchesSnapshot.Branches.HasLocalBranch(branchName)
			hasRemoteBranch := branchesSnapshot.Branches.HasMatchingTrackingBranchFor(branchName, repo.UnvalidatedConfig.NormalConfig.DevRemote)
			if !hasLocalBranch && !hasRemoteBranch {
				return observeData{}, fmt.Errorf(messages.BranchDoesntExist, branchName)
			}
			if hasLocalBranch && !hasRemoteBranch {
				return observeData{}, fmt.Errorf(messages.ObserveBranchIsLocal, branchName)
			}
			branchesToObserve.Add(branchName, branchType)
		}
	}
	return observeData{
		branchInfos:       branchesSnapshot.Branches,
		branchesSnapshot:  branchesSnapshot,
		branchesToObserve: observeCandidates,
		checkout:          branchToCheckout,
	}, err
}
