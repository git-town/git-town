package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/v12/src/cli/dialog"
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
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
			for i := 0; i < int(amount); i++ {
				localBranches = append(localBranches, gitdomain.NewLocalBranchName(fmt.Sprintf("branch-%d", i)))
			}
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{}
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err = dialog.Parent(dialog.ParentArgs{
				Branch:          gitdomain.NewLocalBranchName("branch-2"),
				DialogTestInput: dialogTestInputs.Next(),
				Lineage:         lineage,
				LocalBranches:   localBranches,
				MainBranch:      main,
			})
			return err
		},
	}
}
