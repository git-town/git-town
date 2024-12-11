package debug

import (
	"os"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterCreatePrototypeBranches() *cobra.Command {
	return &cobra.Command{
		Use: "create-prototype-branches",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.NewBranchType(configdomain.BranchTypeFeatureBranch, dialogTestInputs.Next())
			return err
		},
	}
}
