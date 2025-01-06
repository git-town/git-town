package debug

import (
	"os"
	"time"

	"github.com/git-town/git-town/v17/internal/cli/dialog"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/spf13/cobra"
)

func unfinishedStateCommitAuthorCmd() *cobra.Command {
	return &cobra.Command{
		Use: "unfinished-state",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.AskHowToHandleUnfinishedRunState("sync", "feature-branch", time.Now().Add(time.Second*-1), true, dialogTestInputs.Next())
			return err
		},
	}
}
