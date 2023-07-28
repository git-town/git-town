package cmd

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/spf13/cobra"
)

const statusDesc = "Displays or resets the current suspended Git Town command"

func statusCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "status",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   statusDesc,
		Long:    long(statusDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return status(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	cmd.AddCommand(resetRunstateCommand())
	return &cmd
}

func status(debug bool) error {
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 false,
		HandleUnfinishedState: false,
		OmitBranchNames:       false,
		ValidateIsOnline:      false,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	config, err := loadDisplayStatusConfig(repo.RootDir)
	if err != nil {
		return err
	}
	displayStatus(*config)
	repo.Runner.Stats.PrintAnalysis()
	return nil
}

type displayStatusConfig struct {
	filepath string             // filepath of the runstate file
	state    *runstate.RunState // content of the runstate file
}

func loadDisplayStatusConfig(rootDir string) (*displayStatusConfig, error) {
	filepath, err := runstate.PersistenceFilePath(rootDir)
	if err != nil {
		return nil, err
	}
	state, err := runstate.Load(rootDir)
	if err != nil {
		return nil, err
	}
	return &displayStatusConfig{
		filepath: filepath,
		state:    state,
	}, nil
}

func displayStatus(config displayStatusConfig) {
	if config.state == nil {
		fmt.Println("No status file found for this repository.")
		return
	}
	if config.state.IsUnfinished() {
		displayUnfinishedStatus(config)
	} else {
		displayFinishedStatus(config)
	}
}

func displayUnfinishedStatus(config displayStatusConfig) {
	timeDiff := time.Since(config.state.UnfinishedDetails.EndTime)
	fmt.Printf("The last Git Town command (%s) hit a problem %v ago.\n", config.state.Command, timeDiff)
	if config.state.HasAbortSteps() {
		fmt.Println("You can run \"git town abort\" to abort it.")
	}
	if config.state.HasRunSteps() {
		fmt.Println("You can run \"git town continue\" to finish it.")
	}
	if config.state.UnfinishedDetails.CanSkip {
		fmt.Println("You can run \"git town skip\" to skip the currently failing step.")
	}
}

func displayFinishedStatus(config displayStatusConfig) {
	fmt.Printf("The previous Git Town command (%s) finished successfully.\n", config.state.Command)
	if config.state.HasUndoSteps() {
		fmt.Println("You can run \"git town undo\" to undo it.")
	}
}
