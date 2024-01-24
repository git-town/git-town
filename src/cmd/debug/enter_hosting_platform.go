package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/spf13/cobra"
)

func enterHostingPlatform() *cobra.Command {
	return &cobra.Command{
		Use: "hosting-platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := dialog.EnterHostingPlatform("", dialogInputs.Next())
			return err
		},
	}
}
