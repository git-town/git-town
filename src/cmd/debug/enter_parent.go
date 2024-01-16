package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

func enterParentCmd() *cobra.Command {
	return &cobra.Command{
		Use: "enter-parent",
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
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{}
			localBranches := gitdomain.LocalBranchNames{branch1, branch2, branch3, branch4, branch5, branch6, branch7, branch8, branch9, branchA}
			dialogTestInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := dialog.EnterParent(dialog.EnterParentArgs{
				Branch:          branch2,
				DialogTestInput: dialogTestInputs.Next(),
				LocalBranches:   localBranches,
				Lineage:         lineage,
				MainBranch:      main,
			})
			return err
		},
	}
}
