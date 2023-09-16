package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/spf13/cobra"
)

const skipDesc = "Restarts the last run git-town command by skipping the current branch"

func skipCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "skip",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   skipDesc,
		Long:    long(skipDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return skip(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func skip(debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	lineage := repo.Runner.Config.Lineage()
	_, initialBranchesSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  &repo,
		Fetch:                 false,
		HandleUnfinishedState: false,
		Lineage:               lineage,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	runState, err := runstate.Load(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return fmt.Errorf(messages.SkipNothingToDo)
	}
	if !runState.UnfinishedDetails.CanSkip {
		return fmt.Errorf(messages.SkipBranchHasConflicts)
	}
	skipRunState := runState.CreateSkipRunState()
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:                &skipRunState,
		Run:                     &repo.Runner,
		Connector:               nil,
		Lineage:                 lineage,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
	})
}
