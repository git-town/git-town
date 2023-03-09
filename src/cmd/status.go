package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func statusCommand(repo *git.ProdRepo) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, hasGitVersion, isRepository),
		Short:   "Displays or resets the current suspended Git Town command",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadDisplayStatusConfig(repo)
			if err != nil {
				return err
			}
			displayStatus(*config)
			return nil
		},
	}
	cmd.AddCommand(resetRunstateCommand(repo))
	return cmd
}

type displayStatusConfig struct {
	filepath string             // filepath of the runstate file
	state    *runstate.RunState // content of the runstate file
}

func loadDisplayStatusConfig(repo *git.ProdRepo) (*displayStatusConfig, error) {
	filepath, err := runstate.PersistenceFilePath(repo)
	if err != nil {
		return nil, fmt.Errorf("cannot determine the runstate file path: %w", err)
	}
	state, err := runstate.Load(repo)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("the runstate file contains invalid content: %w", err)
		}
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
