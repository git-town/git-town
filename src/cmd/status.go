package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func statusCommand(repo *git.ProdRepo) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Displays or resets the current suspended Git Town command",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := loadDisplayStatusConfig(repo)
			if err != nil {
				cli.Exit(err)
			}
			displayStatus(*config)
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
		GroupID: "errors",
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
	fmt.Printf("The status for this repository is at %s.\n", config.filepath)
	if config.state == nil {
		fmt.Println("No status file found for this repository.")
		return
	}
	fmt.Printf("The previous Git Town command (%s) ", config.state.Command)
	if config.state.IsUnfinished() {
		fmt.Println("did not finish.")
	} else {
		fmt.Println("finished successfully.")
	}
	if config.state.HasAbortSteps() {
		fmt.Println("You can run \"git town abort\" to abort it.")
	}
	if config.state.HasRunSteps() {
		fmt.Println("You can run \"git town continue\" to finish it.")
	}
	if config.state.UnfinishedDetails != nil && config.state.UnfinishedDetails.CanSkip {
		fmt.Println("You can run \"git town skip\" to skip the currently failing step.")
	}
	if config.state.HasUndoSteps() {
		fmt.Println("You can run \"git town undo\" to undo it.")
	}
}
