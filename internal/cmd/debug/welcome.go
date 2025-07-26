package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/spf13/cobra"
)

func welcome() *cobra.Command {
	return &cobra.Command{
		Use: "welcome",
		RunE: func(_ *cobra.Command, _ []string) error {
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, err := dialog.Welcome(inputs)
			return err
		},
	}
}
