package debug

import (
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/spf13/cobra"
)

func enterGiteaToken() *cobra.Command {
	return &cobra.Command{
		Use: "gitea-token",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.GiteaToken(gohacks.NewOptionNone[configdomain.GiteaToken](), dialogInputs.Next())
			return err
		},
	}
}
