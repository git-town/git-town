package cmd

import (
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	configInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/config"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const mergeDesc = "Merges the current branch in a stack with its parent"

const mergeHelp = `
Merges the current branch with its parent branch.
Both branches must be feature branches and be in sync with their tracking branches.
`

func mergeCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "merge",
		Args:    cobra.NoArgs,
		GroupID: "refactoring",
		Short:   mergeDesc,
		Long:    cmdhelpers.Long(mergeDesc, mergeHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeMerge(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeMerge(args []string, verbose configdomain.Verbose) error {
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
	data, err := determineMergeData(args, repo)
	if err != nil {
		return err
	}
	err = validateMergeData(data)
	if err != nil {
		return err
	}
	branchNames := data.branchesToObserve.Keys()
	if err = repo.UnvalidatedConfig.NormalConfig.AddToObservedBranches(branchNames...); err != nil {
		return err
	}
	if err = removeNonObserveBranchTypes(data.branchesToObserve, repo.UnvalidatedConfig); err != nil {
		return err
	}
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

type mergeData struct {
	branchInfos       gitdomain.BranchInfos
	branchesSnapshot  gitdomain.BranchesSnapshot
	branchesToObserve configdomain.BranchesAndTypes
	checkout          Option[gitdomain.LocalBranchName]
}

func determineMergeData(args []string, repo execute.OpenRepoResult) (mergeData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return mergeData{}, err
	}
	branchesToObserve, branchToCheckout, err := execute.BranchesToMark(args, branchesSnapshot, repo.UnvalidatedConfig)
	return mergeData{
		branchInfos:       branchesSnapshot.Branches,
		branchesSnapshot:  branchesSnapshot,
		branchesToObserve: branchesToObserve,
		checkout:          branchToCheckout,
	}, err
}

func validateMergeData(data mergeData) error {
	return nil
}
