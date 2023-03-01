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

func runstateCommand(repo *git.ProdRepo) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Displays or resets the current interrupted Git Town command",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := loadDisplayRunstateConfig(repo)
			if err != nil {
				cli.Exit(err)
			}
			displayRunstate(*config)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
	}
	cmd.AddCommand(resetRunstateCommand(repo))
	return cmd
}

type displayRunstateConfig struct {
	filepath  string
	persisted *runstate.RunState
}

func loadDisplayRunstateConfig(repo *git.ProdRepo) (*displayRunstateConfig, error) {
	filepath, err := runstate.PersistenceFilename(repo)
	if err != nil {
		return nil, fmt.Errorf("cannot determine the runstate file: %w", err)
	}
	persisted, err := runstate.Load(repo)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("the runstate file contains invalid content: %w", err)
		}
	}
	return &displayRunstateConfig{
		filepath:  filepath,
		persisted: persisted,
	}, nil
}

func displayRunstate(config displayRunstateConfig) {
	fmt.Printf("The status for this repository is stored in %s.\n", config.filepath)
	if config.persisted == nil {
		fmt.Println("No status found for this repository.")
		return
	}
	fmt.Printf("The previous Git Town command (%s) ", config.persisted.Command)
	if config.persisted.IsUnfinished() {
		fmt.Println("did not finish.")
	} else {
		fmt.Println("finished successfully.")
	}
	if config.persisted.HasAbortSteps() {
		fmt.Println("You can run \"git town abort\" to abort it.")
	}
	if config.persisted.HasRunSteps() {
		fmt.Println("You can run \"git town continue\" to finish it.")
	}
	if config.persisted.HasUndoSteps() {
		fmt.Println("You can run \"git town undo\" to undo it.")
	}
}
