package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/spf13/cobra"
)

func enterGitHubToken() *cobra.Command {
	return &cobra.Command{
		Use: "github-token",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogInputs := dialog.LoadTestInputs(os.Environ())
			_, err := dialog.EnterGitHubToken(dialogInputs.Next())
			return err
		},
	}
}
