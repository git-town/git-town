package debug

import (
	"os"

	"github.com/git-town/git-town/v20/internal/cli/dialog"
	"github.com/git-town/git-town/v20/internal/cli/dialog/components"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterUnknownBranch() *cobra.Command {
	return &cobra.Command{
		Use: "unknown-branch-type",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.UnknownBranchType(configdomain.BranchTypeFeatureBranch, dialogTestInputs.Next())
			return err
		},
	}
}
