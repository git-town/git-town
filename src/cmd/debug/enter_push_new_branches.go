package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v11/src/cli/dialog/dialogscreens"
	"github.com/spf13/cobra"
)

func enterPushNewBranches() *cobra.Command {
	return &cobra.Command{
		Use: "push-new-branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialogscreens.EnterPushNewBranches(true, dialogTestInputs.Next())
			return err
		},
	}
}
