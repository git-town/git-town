package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/commandconfig"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/config"
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
	if !data.checkout.IsEmpty() {
		if err = repo.Frontend.CheckoutBranch(data.checkout, false); err != nil {
			return err
		}
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:             repo.Backend,
		BeginConfigSnapshot: repo.ConfigSnapshot,
		Command:             "observe",
		CommandsCounter:     repo.CommandsCounter,
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		FinalMessages:       repo.FinalMessages,
		RootDir:             repo.RootDir,
		Verbose:             verbose,
	})
}

type observeData struct {
	allBranches       gitdomain.BranchInfos
	branchesToObserve commandconfig.BranchesAndTypes
	checkout          gitdomain.LocalBranchName
}

func printObservedBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.ObservedBranchIsNowObserved, branch)
	}
}

func removeNonObserveBranchTypes(branches map[gitdomain.LocalBranchName]configdomain.BranchType, config config.UnvalidatedConfig) error {
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

func determineObserveData(args []string, repo execute.OpenRepoResult) (observeData, error) {
	branchesSnapshot, err := repo.Backend.BranchesSnapshot()
	if err != nil {
		return observeData{}, err
	}
	branchesToObserve := commandconfig.BranchesAndTypes{}
	checkout := gitdomain.EmptyLocalBranchName()
	switch len(args) {
	case 0:
		branchesToObserve.Add(branchesSnapshot.Active, *repo.UnvalidatedConfig.Config)
	case 1:
		branch := gitdomain.NewLocalBranchName(args[0])
		branchesToObserve.Add(branch, *repo.UnvalidatedConfig.Config)
		branchInfo := branchesSnapshot.Branches.FindByRemoteName(branch.TrackingBranch())
		if branchInfo.SyncStatus == gitdomain.SyncStatusRemoteOnly {
			checkout = branch
		}
	default:
		branchesToObserve.AddMany(gitdomain.NewLocalBranchNames(args...), *repo.UnvalidatedConfig.Config)
	}
	return observeData{
		allBranches:       branchesSnapshot.Branches,
		branchesToObserve: branchesToObserve,
		checkout:          checkout,
	}, nil
}

func validateObserveData(data observeData) error {
	for branchName, branchType := range data.branchesToObserve {
		if !data.allBranches.HasLocalBranch(branchName) && !data.allBranches.HasMatchingTrackingBranchFor(branchName) {
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
