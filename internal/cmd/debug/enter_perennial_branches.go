package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/spf13/cobra"
)

func enterPerennialBranches() *cobra.Command {
	return &cobra.Command{
		Use:  "perennial-branches <number of branches>",
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
			existingPerennialBranches := gitdomain.NewLocalBranchNames("branch-2", "branch-4")
			main := gitdomain.NewLocalBranchName("main")
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err = dialog.PerennialBranches(localBranches, existingPerennialBranches, main, dialogTestInputs.Next())
			return err
		},
	}
}
