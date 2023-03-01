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
		Use:   "runstate",
		Short: "Displays or resets the persisted runstate",
		Run: func(cmd *cobra.Command, args []string) {
			filepath, err := runstate.PersistenceFilename(repo)
			if err != nil {
				cli.Exit(fmt.Errorf("cannot determine the runstate file: %w", err))
			}
			fmt.Printf("The runstate is stored in %s.\n", filepath)
			_, err = os.Stat(filepath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					fmt.Println("No runstate exists.")
					return
				} else {
					cli.Exit(fmt.Errorf("cannot analyze runstate: %w", err))
				}
			}
			fmt.Println("The runstate file exists.")
			_, err = runstate.Load(repo)
			if err != nil {
				fmt.Printf("Cannot load current runstate: %v\n", err)
			} else {
				fmt.Println("Runstate is valid.")
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
	}
	cmd.AddCommand(resetRunstateCommand(repo))
	return cmd
}
