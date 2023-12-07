package cmd

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/git-town/git-town/v11/src/vm/statefile"
	"github.com/spf13/cobra"
)

const statusDesc = "Displays or resets the current suspended Git Town command"

func statusCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "status",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   statusDesc,
		Long:    long(statusDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeStatus(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	cmd.AddCommand(resetRunstateCommand())
	return &cmd
}

func executeStatus(verbose bool) error {
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
	config, err := loadDisplayStatusConfig(repo.RootDir)
	if err != nil {
		return err
	}
	displayStatus(*config)
	print.Footer(verbose, repo.Runner.CommandsCounter.Count(), print.NoFinalMessages)
	return nil
}

type displayStatusConfig struct {
	filepath string             // filepath of the runstate file
	state    *runstate.RunState // content of the runstate file
}

func loadDisplayStatusConfig(rootDir domain.RepoRootDir) (*displayStatusConfig, error) {
	filepath, err := statefile.FilePath(rootDir)
	if err != nil {
		return nil, err
	}
	state, err := statefile.Load(rootDir)
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
	if config.state.HasAbortProgram() {
		fmt.Println("You can run \"git town undo\" to go back to where you started.")
	}
	if config.state.HasRunProgram() {
		fmt.Println("You can run \"git town continue\" to finish it.")
	}
	if config.state.UnfinishedDetails.CanSkip {
		fmt.Println("You can run \"git town skip\" to skip the currently failing operation.")
	}
}

func displayFinishedStatus(config displayStatusConfig) {
	fmt.Printf("The previous Git Town command (%s) finished successfully.\n", config.state.Command)
	if config.state.HasUndoProgram() {
		fmt.Println("You can run \"git town undo\" to undo it.")
	}
}
