package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
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
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{}
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err = dialog.Parent(dialog.ParentArgs{
				Branch:          gitdomain.NewLocalBranchName("branch-2"),
				DefaultChoice:   main,
				DialogTestInput: dialogTestInputs.Value.Next(),
				Lineage:         lineage,
				LocalBranches:   localBranches,
				MainBranch:      main,
			})
			return err
		},
	}
}
