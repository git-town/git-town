package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/spf13/cobra"
)

func enterParentCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "parent <number of branches>",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			amount, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			localBranches := gitdomain.LocalBranchNames{}
			for i := range amount {
				localBranches = append(localBranches, gitdomain.NewLocalBranchName(fmt.Sprintf("branch-%d", i)))
			}
			lineage := configdomain.Lineage{}
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err = dialog.Parent(dialog.ParentArgs{
				Branch:          "branch-2",
				DefaultChoice:   "main",
				DialogTestInput: dialogTestInputs.Next(),
				Lineage:         lineage,
				LocalBranches:   localBranches,
				MainBranch:      "main",
			})
			return err
		},
	}
}
