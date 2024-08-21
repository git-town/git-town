package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/flags"
	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v15/internal/config"
	"github.com/git-town/git-town/v15/internal/config/commandconfig"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/execute"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	configInterpreter "github.com/git-town/git-town/v15/internal/vm/interpreter/config"
	. "github.com/git-town/git-town/v15/pkg/prelude"
	"github.com/spf13/cobra"
)

const observeDesc = "Stop your contributions to some feature branches"

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

func executeObserve(args []string, verbose configdomain.Verbose) error {
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
	data, err := determineObserveData(args, repo)
	if err != nil {
		return err
	}
	err = validateObserveData(data)
	if err != nil {
		return err
	}
	branchNames := data.branchesToObserve.Keys()
	if err = repo.UnvalidatedConfig.AddToObservedBranches(branchNames...); err != nil {
		return err
	}
	if err = removeNonObserveBranchTypes(data.branchesToObserve, repo.UnvalidatedConfig); err != nil {
		return err
	}
	printObservedBranches(branchNames)
	if checkout, hasCheckout := data.checkout.Get(); hasCheckout {
		if err = repo.Git.CheckoutBranch(repo.Frontend, checkout, false); err != nil {
			return err
		}
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: Some(data.branchesSnapshot),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "observe",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       branchNames.BranchNames(),
		Verbose:               verbose,
	})
}

type observeData struct {
	allBranches       gitdomain.BranchInfos
	branchesSnapshot  gitdomain.BranchesSnapshot
	branchesToObserve commandconfig.BranchesAndTypes
	checkout          Option[gitdomain.LocalBranchName]
}

func printObservedBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.ObservedBranchIsNowObserved, branch)
	}
}

func removeNonObserveBranchTypes(branches commandconfig.BranchesAndTypes, config config.UnvalidatedConfig) error {
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
		case configdomain.BranchTypePrototypeBranch:
			if err := config.RemoveFromPrototypeBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		}
	}
	return nil
}

func determineObserveData(args []string, repo execute.OpenRepoResult) (observeData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return observeData{}, err
	}
	branchesToObserve := commandconfig.BranchesAndTypes{}
	branchToCheckout := None[gitdomain.LocalBranchName]()
	switch len(args) {
	case 0:
		currentBranch, hasCurrentBranch := branchesSnapshot.Active.Get()
		if !hasCurrentBranch {
			return observeData{}, errors.New(messages.CurrentBranchCannotDetermine)
		}
		branchesToObserve.Add(currentBranch, repo.UnvalidatedConfig.Config.Get())
	case 1:
		branch := gitdomain.NewLocalBranchName(args[0])
		branchesToObserve.Add(branch, repo.UnvalidatedConfig.Config.Get())
		branchInfo, hasBranchInfo := branchesSnapshot.Branches.FindByRemoteName(branch.TrackingBranch()).Get()
		if hasBranchInfo && branchInfo.SyncStatus == gitdomain.SyncStatusRemoteOnly {
			branchToCheckout = Some(branch)
		} else {
			branchToCheckout = None[gitdomain.LocalBranchName]()
		}
	default:
		branchesToObserve.AddMany(gitdomain.NewLocalBranchNames(args...), repo.UnvalidatedConfig.Config.Get())
	}
	return observeData{
		allBranches:       branchesSnapshot.Branches,
		branchesSnapshot:  branchesSnapshot,
		branchesToObserve: branchesToObserve,
		checkout:          branchToCheckout,
	}, nil
}

func validateObserveData(data observeData) error {
	for branchName, branchType := range data.branchesToObserve {
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotObserve)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotObserve)
		case configdomain.BranchTypeObservedBranch:
			return fmt.Errorf(messages.BranchIsAlreadyObserved, branchName)
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeContributionBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
		}
		hasLocalBranch := data.allBranches.HasLocalBranch(branchName)
		hasRemoteBranch := data.allBranches.HasMatchingTrackingBranchFor(branchName)
		if !hasLocalBranch && !hasRemoteBranch {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
		if hasLocalBranch && !hasRemoteBranch {
			return fmt.Errorf(messages.ObserveBranchIsLocal, branchName)
		}
	}
	return nil
}
