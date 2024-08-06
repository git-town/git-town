package debug

import (
	"os"
	"time"

	"github.com/git-town/git-town/v14/internal/cli/dialog"
	"github.com/git-town/git-town/v14/internal/cli/dialog/components"
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/spf13/cobra"
)

func unfinishedStateCommitAuthorCmd() *cobra.Command {
	return &cobra.Command{
		Use: "unfinished-state",
		RunE: func(_ *cobra.Command, _ []string) error {
			branch := gitdomain.NewLocalBranchName("feature-branch")
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.AskHowToHandleUnfinishedRunState("sync", branch, time.Now().Add(time.Second*-1), true, dialogTestInputs.Next())
			return err
		},
	}
}
