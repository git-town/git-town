package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/cli/dialogs/enter"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/spf13/cobra"
)

func enterGiteaToken() *cobra.Command {
	return &cobra.Command{
		Use: "gitea-token",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := enter.GiteaToken(configdomain.GiteaToken(""), dialogInputs.Next())
			return err
		},
	}
}
