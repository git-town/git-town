package debug

import (
	"os"

	"github.com/git-town/git-town/v14/internal/cli/dialog"
	"github.com/git-town/git-town/v14/internal/cli/dialog/components"
	"github.com/spf13/cobra"
)

func enterCreatePrototypeBranches() *cobra.Command {
	return &cobra.Command{
		Use: "create-prototype-branches",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.CreatePrototypeBranches(true, dialogTestInputs.Next())
			return err
		},
	}
}
