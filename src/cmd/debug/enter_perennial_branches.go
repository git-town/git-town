package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

func enterPerennialBranchesCmd() *cobra.Command {
	return &cobra.Command{
		Use: "perennial-branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			branch3 := gitdomain.NewLocalBranchName("branch-3")
			branch4 := gitdomain.NewLocalBranchName("branch-4")
			branch5 := gitdomain.NewLocalBranchName("branch-5")
			branch6 := gitdomain.NewLocalBranchName("branch-6")
			branch7 := gitdomain.NewLocalBranchName("branch-7")
			branch8 := gitdomain.NewLocalBranchName("branch-8")
			branch9 := gitdomain.NewLocalBranchName("branch-9")
			branchA := gitdomain.NewLocalBranchName("branch-A")
			localBranches := gitdomain.LocalBranchNames{branch1, branch2, branch3, branch4, branch5, branch6, branch7, branch8, branch9, branchA}
			existingPerennialBranches := gitdomain.LocalBranchNames{branch1, branch4}
			main := gitdomain.NewLocalBranchName("main")
			dialogTestInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := dialog.EnterPerennialBranches(localBranches, existingPerennialBranches, main, dialogTestInputs.Next())
			return err
		},
	}
}
