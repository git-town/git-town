package debug

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

func selectCommitAuthorCmd() *cobra.Command {
	return &cobra.Command{
		Use: "select-commit-author",
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := gitdomain.NewLocalBranchName("feature-branch")
			authors := []string{"Jean-Luc Picard <captain@enterprise.com>", "William Riker <numberone@enterprise.com>"}
			selected, aborted, err := dialog.SelectSquashCommitAuthor(branch, authors)
			if err != nil {
				return err
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
