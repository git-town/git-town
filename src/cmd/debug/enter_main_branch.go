package debug

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

func enterMainBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use: "enter-main-branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			localBranches := gitdomain.NewLocalBranchNames("branch-1", "branch-2", "branch-3", "branch-4", "branch-5", "branch-6", "branch-7", "branch-8", "branch-9", "branch-A", "branch-B", "main")
			main := gitdomain.NewLocalBranchName("main")
			selected, aborted, err := dialog.EnterMainBranch(localBranches, main)
			if err != nil {
				return err
			}
			if aborted {
				fmt.Println("ABORTED")
			}
			if aborted {
				fmt.Println("ABORTED")
			} else {
				fmt.Println("SELECTED:", selected)
			}
			return nil
		},
	}
}
