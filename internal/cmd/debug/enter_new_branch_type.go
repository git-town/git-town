package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterNewBranchType() *cobra.Command {
	return &cobra.Command{
		Use: "new-branch-type",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.NewBranchType(dialog.Args[configdomain.NewBranchType]{
				Global: Some(configdomain.NewBranchType(configdomain.BranchTypeFeatureBranch)),
				Inputs: dialogTestInputs,
				Local:  Some(configdomain.NewBranchType(configdomain.BranchTypePrototypeBranch)),
			})
			return err
		},
	}
}
