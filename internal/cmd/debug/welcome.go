package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/spf13/cobra"
)

func welcome() *cobra.Command {
	return &cobra.Command{
		Use: "welcome",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, err := dialog.Welcome(dialogTestInputs.Next())
			return err
		},
	}
}
