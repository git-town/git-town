package cmd

import (
	"errors"
	"fmt"

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
			isRepo := repo.Silent.IsRepository()
			if !isRepo {
				cli.Exit(errors.New("running outside a Git repository"))
			}
			cli.Print("Hosting service: ")
			connector, err := hosting.NewConnector(&repo.Config, &repo.Silent, cli.PrintConnectorAction)
			if err != nil {
				cli.Exit(fmt.Errorf("(cannot determine hosting connector: %v)", err))
			} else {
				cli.Println(connector.HostingServiceName())
			}
		},
		Args: cobra.NoArgs,
	}
}
