package debug

import (
	"os"
	"time"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/spf13/cobra"
)

func unfinishedStateCommitAuthorCmd() *cobra.Command {
	return &cobra.Command{
		Use: "unfinished-state",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.AskHowToHandleUnfinishedRunState("sync", "feature-branch", time.Now().Add(time.Second*-1), true, dialogTestInputs.Next())
			return err
		},
	}
}
