package cmd

import (
	"errors"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/spf13/cobra"
)

const setParentDesc = "Prompts to set the parent branch for the current branch"

func setParentCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "set-parent",
		GroupID: "lineage",
		Args:    cobra.NoArgs,
		Short:   setParentDesc,
		Long:    cmdhelpers.Long(setParentDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSetParent(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSetParent(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	branches, _, _, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	if !repo.Runner.IsFeatureBranch(branches.Initial) {
		return errors.New(messages.SetParentNoFeatureBranch)
	}
	existingParent := repo.Runner.Lineage.Parent(branches.Initial)
	if !existingParent.IsEmpty() {
		// TODO: delete the old parent only when the user has entered a new parent
		repo.Runner.Config.RemoveParent(branches.Initial)
		repo.Runner.Config.Reload()
	} else {
		existingParent = repo.Runner.Config.MainBranch
	}
	err = execute.EnsureKnownBranchAncestry(branches.Initial, execute.EnsureKnownBranchAncestryArgs{
		Config:        &repo.Runner.FullConfig,
		AllBranches:   branches.All,
		DefaultBranch: existingParent,
		Runner:        repo.Runner,
	})
	if err != nil {
		return err
	}
	print.Footer(verbose, repo.Runner.CommandsCounter.Count(), print.NoFinalMessages)
	return nil
}
