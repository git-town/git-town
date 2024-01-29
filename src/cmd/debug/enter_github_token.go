package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/spf13/cobra"
)

func enterGitHubToken() *cobra.Command {
	return &cobra.Command{
		Use: "github-token",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialogs.GitHubToken(configdomain.GitHubToken(""), dialogInputs.Next())
			return err
		},
	}
}
