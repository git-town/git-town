package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterShareNewBranches() *cobra.Command {
	return &cobra.Command{
		Use: "share-new-branches",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.ShareNewBranches(configdomain.ShareNewBranchesNone, dialogTestInputs)
			return err
		},
	}
}
