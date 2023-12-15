package cmd

import (
	"errors"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/print"
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
		Long:    long(setParentDesc),
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
	lineage := repo.Runner.GitTown.Lineage(repo.Runner.Backend.GitTown.RemoveLocalConfigValue)
	pushHook, err := repo.Runner.GitTown.PushHook()
	if err != nil {
		return err
	}
	branches, _, _, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	if !branches.Types.IsFeatureBranch(branches.Initial) {
		return errors.New(messages.SetParentNoFeatureBranch)
	}
	existingParent := lineage.Parent(branches.Initial)
	if !existingParent.IsEmpty() {
		// TODO: delete the old parent only when the user has entered a new parent
		repo.Runner.GitTown.RemoveParent(branches.Initial)
	} else {
		existingParent = repo.Runner.GitTown.MainBranch()
	}
	mainBranch := repo.Runner.GitTown.MainBranch()
	branches.Types, _, err = execute.EnsureKnownBranchAncestry(branches.Initial, execute.EnsureKnownBranchAncestryArgs{
		AllBranches:   branches.All,
		BranchTypes:   branches.Types,
		DefaultBranch: existingParent,
		Lineage:       lineage,
		MainBranch:    mainBranch,
		Runner:        repo.Runner,
	})
	if err != nil {
		return err
	}
	print.Footer(verbose, repo.Runner.CommandsCounter.Count(), print.NoFinalMessages)
	return nil
}
