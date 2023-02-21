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
		Short: "Analyzes the Git Town setup",
		Run: func(cmd *cobra.Command, args []string) {
			connector, err := hosting.NewConnector(&repo.Config, &repo.Silent, cli.PrintDriverAction)
			if err != nil {
				cli.Exit(err)
			}
			cli.Printf("Git Town v%s\n", version)
			cli.Println()
			cli.PrintHeader("code hosting platform")
			cli.PrintEntry("address", connector.RepositoryURL())
			cli.PrintEntry("type", connector.HostingServiceName())
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			if err := validateIsConfigured(repo); err != nil {
				return err
			}
			return repo.Config.ValidateIsOnline()
		},
	}
}
