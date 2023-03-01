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

func resetRunstateCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Resets the current suspended Git Town command",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := loadResetStatusConfig(repo)
			if err != nil {
				cli.Exit(err)
			}
			err = resetStatus(config)
			if err != nil {
				cli.Exit(fmt.Errorf("cannot delete runstate file: %w", err))
			}
			fmt.Println("Runstate file deleted.")
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
	}
}

type resetStatusConfig struct {
	filepath string // filepath of the runstate file
}

func loadResetStatusConfig(repo *git.ProdRepo) (*resetStatusConfig, error) {
	filepath, err := runstate.PersistenceFilename(repo)
	if err != nil {
		return nil, fmt.Errorf("cannot determine the runstate file path: %w", err)
	}
	return &resetStatusConfig{
		filepath: filepath,
	}, nil
}

func resetStatus(config *resetStatusConfig) error {
	err := os.Remove(config.filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Runstate doesn't exist.")
			return nil
		}
		return fmt.Errorf("cannot delete runstate file: %w", err)
	}
	return nil
}
