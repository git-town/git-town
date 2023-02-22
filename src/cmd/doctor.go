package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/spf13/cobra"
)

func doctorCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Displays diagnostic information",
		Run: func(cmd *cobra.Command, args []string) {
			cli.Printf("Git Town v%s\n", version)
			cli.Println()
			err := ValidateIsRepository(repo)
			if err != nil {
				return err
			}
			if err := validateIsConfigured(repo); err != nil {
				return err
			}
			if err := repo.Config.ValidateIsOnline(); err != nil {
				return err
			}
			driver, driverErr := hosting.NewDriver(&repo.Config, &repo.Silent, cli.PrintDriverAction)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.NoArgs,
	}
}
