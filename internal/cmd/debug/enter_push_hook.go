package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/spf13/cobra"
)

func enterPushHookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "push-hook",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.PushHook(true, dialogTestInputs.Next())
			return err
		},
	}
}
