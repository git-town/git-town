package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterMainBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use: "main-branch",
		RunE: func(_ *cobra.Command, _ []string) error {
			localBranches := gitdomain.NewLocalBranchNames("main", "branch-1", "branch-2", "branch-3", "branch-4", "branch-5", "branch-6", "branch-7", "branch-8", "branch-9", "branch-A", "branch-B")
			main := gitdomain.NewLocalBranchName("main")
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, _, _, err := dialog.MainBranch(dialog.MainBranchArgs{
				Inputs:         inputs,
				Local:          Some(main),
				LocalBranches:  localBranches,
				StandardBranch: Some(gitdomain.NewLocalBranchName("main")),
				Unscoped:       None[gitdomain.LocalBranchName](),
			})
			return err
		},
	}
}
