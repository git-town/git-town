package debug

import (
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/spf13/cobra"
)

func enterPushNewBranches() *cobra.Command {
	return &cobra.Command{
		Use: "push-new-branches",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.PushNewBranches(true, dialogTestInputs.Value.Next())
			return err
		},
	}
}
