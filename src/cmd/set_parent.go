package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/spf13/cobra"
)

const setParentDesc = "Prompt to set the parent branch for the current branch"

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
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Runner.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
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
	existingParentPtr := repo.Runner.Config.FullConfig.Lineage.Parent(branchesSnapshot.Active)
	var defaultBranch gitdomain.LocalBranchName
	if existingParentPtr != nil {
		defaultBranch = *existingParentPtr
		// TODO: delete the old parent only when the user has entered a new parent
		repo.Runner.Config.RemoveParent(branchesSnapshot.Active)
		repo.Runner.Config.Reload()
	} else {
		defaultBranch = repo.Runner.Config.FullConfig.MainBranch
	}
	err = execute.EnsureKnownBranchAncestry(branchesSnapshot.Active, execute.EnsureKnownBranchAncestryArgs{
		AllBranches:      branchesSnapshot.Branches,
		Config:           repo.Runner.Config,
		DefaultChoice:    existingParent,
		DialogTestInputs: &dialogTestInputs,
		MainBranch:       repo.Runner.Config.FullConfig.MainBranch,
		Runner:           repo.Runner,
	})
	if err != nil {
		return err
	}
	print.Footer(verbose, repo.Runner.CommandsCounter.Count(), print.NoFinalMessages)
	return nil
}
