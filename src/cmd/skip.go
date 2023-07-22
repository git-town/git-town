package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
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
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:                debug,
		DryRun:               false,
		OmitBranchNames:      false,
		ValidateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	_, _, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsOnline:      false,
	})
	if err != nil || exit {
		return err
	}
	runState, err := runstate.Load(&run.Backend)
	if err != nil {
		return fmt.Errorf("cannot load previous run state: %w", err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return fmt.Errorf("nothing to skip")
	}
	if !runState.UnfinishedDetails.CanSkip {
		return fmt.Errorf("cannot skip branch that resulted in conflicts")
	}
	skipRunState := runState.CreateSkipRunState()
	return runstate.Execute(&skipRunState, &run, nil)
}
