package debug

import (
	"fmt"
	"os"
	"time"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

func unfinishedStateCommitAuthorCmd() *cobra.Command {
	return &cobra.Command{
		Use: "unfinished-state",
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := gitdomain.NewLocalBranchName("feature-branch")
			dialogTestInputs := dialog.LoadTestInputs(os.Environ())
			selected, aborted, err := dialog.AskHowToHandleUnfinishedRunState("sync", branch, time.Now().Add(time.Second*-1), true, dialogTestInputs.Next())
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
