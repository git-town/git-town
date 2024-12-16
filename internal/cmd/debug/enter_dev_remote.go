package debug

import (
	"os"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/spf13/cobra"
)

func enterDevRemote() *cobra.Command {
	return &cobra.Command{
		Use: "dev-remote",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.DevRemote("origin", gitdomain.Remotes{"origin", "fork"}, dialogInputs.Next())
			return err
		},
	}
}
