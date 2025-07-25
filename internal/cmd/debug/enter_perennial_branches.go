package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
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
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, _, err = dialog.PerennialBranches(dialog.PerennialBranchesArgs{
				ImmutableGitPerennials: gitdomain.NewLocalBranchNames("global-1", "global-2"),
				Inputs:                 inputs,
				LocalBranches:          localBranches,
				LocalGitPerennials:     gitdomain.NewLocalBranchNames("local-1", "local-2"),
				MainBranch:             "main",
			})
			return err
		},
	}
}
