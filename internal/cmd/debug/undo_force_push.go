package debug

import (
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/spf13/cobra"
)

func undoForcePush() *cobra.Command {
	return &cobra.Command{
		Use: "undo-force-push",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.UndoForcePush("origin/feature", "111111", "222222", dialogTestInputs.Next())
			return err
		},
	}
}
