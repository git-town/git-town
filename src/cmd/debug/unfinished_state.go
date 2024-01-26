package debug

import (
	"os"
	"time"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialogcomponents"
	"github.com/git-town/git-town/v11/src/cli/dialogs/enter"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

func unfinishedStateCommitAuthorCmd() *cobra.Command {
	return &cobra.Command{
		Use: "unfinished-state",
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := gitdomain.NewLocalBranchName("feature-branch")
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := enter.AskHowToHandleUnfinishedRunState("sync", branch, time.Now().Add(time.Second*-1), true, dialogTestInputs.Next())
			return err
		},
	}
}
