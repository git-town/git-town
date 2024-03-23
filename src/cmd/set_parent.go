package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v13/src/cli/dialog/components"
	"github.com/git-town/git-town/v13/src/cli/flags"
	"github.com/git-town/git-town/v13/src/cli/print"
	"github.com/git-town/git-town/v13/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v13/src/execute"
	"github.com/git-town/git-town/v13/src/messages"
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeSetParent(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSetParent(verbose bool) error {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	branchesSnapshot, _, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		FullConfig:            &repo.Runner.Config.FullConfig,
		HandleUnfinishedState: true,
		Repo:                  repo,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return err
	}
	if repo.Runner.Config.FullConfig.IsMainOrPerennialBranch(branchesSnapshot.Active) {
		return fmt.Errorf(messages.SetParentNoFeatureBranch, branchesSnapshot.Active)
	}
	existingParent := repo.Runner.Config.FullConfig.Lineage.Parent(branchesSnapshot.Active)
	if !existingParent.IsEmpty() {
		// TODO: delete the old parent only when the user has entered a new parent
		repo.Runner.Config.RemoveParent(branchesSnapshot.Active)
		repo.Runner.Config.Reload()
	} else {
		existingParent = repo.Runner.Config.FullConfig.MainBranch
	}
	err = execute.EnsureKnownBranchAncestry(branchesSnapshot.Active, execute.EnsureKnownBranchAncestryArgs{
		Config:           &repo.Runner.Config.FullConfig,
		AllBranches:      branchesSnapshot.Branches,
		DefaultBranch:    existingParent,
		DialogTestInputs: &dialogTestInputs,
		Runner:           repo.Runner,
	})
	if err != nil {
		return err
	}
	print.Footer(verbose, repo.Runner.CommandsCounter.Count(), print.NoFinalMessages)
	return nil
}
