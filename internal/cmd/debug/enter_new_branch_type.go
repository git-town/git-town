package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterNewBranchType() *cobra.Command {
	return &cobra.Command{
		Use: "new-branch-type",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.NewBranchType(Some(configdomain.BranchTypePrototypeBranch), dialogTestInputs.Next())
			return err
		},
	}
}
