package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/cli/dialogs/enter"
	"github.com/spf13/cobra"
)

func enterHostingPlatform() *cobra.Command {
	return &cobra.Command{
		Use: "hosting-platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := enter.HostingPlatform("", dialogInputs.Next())
			return err
		},
	}
}
