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
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/config"
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
	data, err := determineContributeData(args, repo)
	if err != nil {
		return err
	}
	err = validateContributeData(data)
	if err != nil {
		return err
	}
	branchNames := data.branchesToMark.Keys()
	if err = repo.Config.AddToContributionBranches(branchNames...); err != nil {
		return err
	}
	if err = removeNonContributionBranchTypes(data.branchesToMark, repo.Config); err != nil {
		return err
	}
	printContributeBranches(branchNames)
	branchToCheckout, hasBranchToCheckout := data.branchToCheckout.Get()
	if hasBranchToCheckout {
		if err = repo.Frontend.CheckoutBranch(branchToCheckout, false); err != nil {
			return err
		}
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:             repo.Backend,
		BeginConfigSnapshot: repo.ConfigSnapshot,
		Command:             "contribute",
		CommandsCounter:     repo.CommandsCounter,
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		FinalMessages:       repo.FinalMessages,
		RootDir:             repo.RootDir,
		Verbose:             verbose,
	})
}

type contributeData struct {
	allBranches      gitdomain.BranchInfos
	branchToCheckout Option[gitdomain.LocalBranchName]
	branchesToMark   commandconfig.BranchesAndTypes
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

func determineContributeData(args []string, repo *execute.OpenRepoResult) (contributeData, error) {
	branchesSnapshot, err := repo.Backend.BranchesSnapshot()
	if err != nil {
		return contributeData{}, err
	}
	branchesToMark := commandconfig.BranchesAndTypes{}
	var branchToCheckout Option[gitdomain.LocalBranchName]
	switch len(args) {
	case 0:
		branchesToMark.Add(branchesSnapshot.Active, &repo.Config.Config)
		branchToCheckout = None[gitdomain.LocalBranchName]()
	case 1:
		branch := gitdomain.NewLocalBranchName(args[0])
		branchesToMark.Add(branch, &repo.Config.Config)
		branchInfo := branchesSnapshot.Branches.FindByRemoteName(branch.TrackingBranch())
		if branchInfo.SyncStatus == gitdomain.SyncStatusRemoteOnly {
			branchToCheckout = Some(branch)
		} else {
			branchToCheckout = None[gitdomain.LocalBranchName]()
		}
	default:
		branchesToMark.AddMany(gitdomain.NewLocalBranchNames(args...), &repo.Config.Config)
		branchToCheckout = None[gitdomain.LocalBranchName]()
	}
	return contributeData{
		allBranches:      branchesSnapshot.Branches,
		branchToCheckout: branchToCheckout,
		branchesToMark:   branchesToMark,
	}, nil
}

func validateContributeData(data contributeData) error {
	for branchName, branchType := range data.branchesToMark {
		if !data.allBranches.HasLocalBranch(branchName) && !data.allBranches.HasMatchingTrackingBranchFor(branchName) {
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
