package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialog"
	"github.com/git-town/git-town/v11/src/cli/dialogs/enter"
	"github.com/spf13/cobra"
)

func enterPushNewBranches() *cobra.Command {
	return &cobra.Command{
		Use: "push-new-branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := enter.PushNewBranches(true, dialogTestInputs.Next())
			return err
		},
	}
}
