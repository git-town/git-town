package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/skip"
	"github.com/git-town/git-town/v12/src/vm/statefile"
	"github.com/spf13/cobra"
)

const skipDesc = "Restarts the last run git-town command by skipping the current branch"

func skipCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "skip",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   skipDesc,
		Long:    cmdhelpers.Long(skipDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSkip(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSkip(verbose bool) error {
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
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	currentSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return err
	}
	runState, err := statefile.Load(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || runState.IsFinished() {
		return errors.New(messages.SkipNothingToDo)
	}
	if !runState.UnfinishedDetails.CanSkip {
		return errors.New(messages.SkipBranchHasConflicts)
	}
	return skip.Execute(skip.ExecuteArgs{
		CurrentBranch:    currentSnapshot.Active,
		HasOpenChanges:   repoStatus.OpenChanges,
		InitialStashSize: runState.BeforeStashSize,
		RootDir:          repo.RootDir,
		RunState:         runState,
		Runner:           repo.Runner,
		TestInputs:       dialogTestInputs,
		Verbose:          verbose,
	})
}
