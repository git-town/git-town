package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

func selectCommitAuthorCmd() *cobra.Command {
	return &cobra.Command{
		Use: "select-commit-author",
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := gitdomain.NewLocalBranchName("feature-branch")
			authors := []string{"Jean-Luc Picard <captain@enterprise.com>", "William Riker <numberone@enterprise.com>"}
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialogs.SelectSquashCommitAuthor(branch, authors, dialogTestInputs.Next())
			return err
		},
	}
}
