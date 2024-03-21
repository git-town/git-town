package debug

import (
	"os"

	"github.com/git-town/git-town/v12/src/cli/dialog"
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/spf13/cobra"
)

func selectCommitAuthorCmd() *cobra.Command {
	return &cobra.Command{
		Use: "select-commit-author",
		RunE: func(_ *cobra.Command, _ []string) error {
			branch := gitdomain.NewLocalBranchName("feature-branch")
			authors := []string{"Jean-Luc Picard <captain@enterprise.com>", "William Riker <numberone@enterprise.com>"}
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.SelectSquashCommitAuthor(branch, authors, dialogTestInputs.Next())
			return err
		},
	}
}
