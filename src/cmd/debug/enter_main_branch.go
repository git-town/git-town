package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

func enterMainBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use: "main-branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			localBranches := gitdomain.NewLocalBranchNames("main", "branch-1", "branch-2", "branch-3", "branch-4", "branch-5", "branch-6", "branch-7", "branch-8", "branch-9", "branch-A", "branch-B")
			main := gitdomain.NewLocalBranchName("main")
			dialogInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := dialog.EnterMainBranch(localBranches, main, dialogInputs.Next())
			return err
		},
	}
}
