package debug

import (
	"os"

	"github.com/git-town/git-town/v17/internal/cli/dialog"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/spf13/cobra"
)

func welcome() *cobra.Command {
	return &cobra.Command{
		Use: "welcome",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, err := dialog.Welcome(dialogTestInputs.Next())
			return err
		},
	}
}
