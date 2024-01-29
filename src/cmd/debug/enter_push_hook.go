package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/cli/dialogs/enter"
	"github.com/spf13/cobra"
)

func enterPushHookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "push-hook",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := enter.PushHook(true, dialogTestInputs.Next())
			return err
		},
	}
}
