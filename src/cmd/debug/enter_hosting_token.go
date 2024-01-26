package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/spf13/cobra"
)

func enterHostingToken() *cobra.Command {
	return &cobra.Command{
		Use: "hosting-token",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogInputs := dialog.LoadTestInputs(os.Environ())
			_, err := dialog.EnterHostingToken(configdomain.CodeHostingPlatformGitHub, dialogInputs.Next())
			return err
		},
	}
}
