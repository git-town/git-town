package debug

import (
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/spf13/cobra"
)

func undoCreateRemoteBranch() *cobra.Command {
	return &cobra.Command{
		Use: "undo-create-remote-branch",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.UndoCreateRemoteBranch("feature", dialogTestInputs.Next())
			return err
		},
	}
}
