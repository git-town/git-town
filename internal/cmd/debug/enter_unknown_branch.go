package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterUnknownBranch() *cobra.Command {
	return &cobra.Command{
		Use: "unknown-branch-type",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.UnknownBranchType(configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch), dialogTestInputs)
			return err
		},
	}
}
