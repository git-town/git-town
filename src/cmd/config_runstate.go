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
			fmt.Printf("The runstate for this repository is stored in %s.\n", filepath)
			_, err = os.Stat(filepath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					fmt.Println("This file doesn't exist.")
					return
				}
				cli.Exit(fmt.Errorf("cannot analyze runstate: %w", err))
			}
			fmt.Print("This file exists ")
			persisted, err := runstate.Load(repo)
			if err != nil {
				cli.Exit(fmt.Errorf("but contains invalid content: %w", err))
			} else {
				fmt.Println("and is valid.")
			}
			fmt.Printf("The previous Git Town command (%s) ", persisted.Command)
			if persisted.IsUnfinished() {
				fmt.Println("did not finish.")
			} else {
				fmt.Println("finished successfully.")
			}
			if persisted.HasAbortSteps() {
				fmt.Println("You can run \"git town abort\" to abort it.")
			}
			if persisted.HasRunSteps() {
				fmt.Println("You can run \"git town continue\" to finish it.")
			}
			if persisted.HasUndoSteps() {
				fmt.Println("You can run \"git town undo\" to undo it.")
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
