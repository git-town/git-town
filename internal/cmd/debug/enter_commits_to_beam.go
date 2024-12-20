package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/v17/internal/cli/dialog"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/spf13/cobra"
)

func enterCommitsToBeam() *cobra.Command {
	return &cobra.Command{
		Use:  "commits-to-beam <number of commits>",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			amount, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			commits := []gitdomain.Commit{}
			for i := range amount {
				commits = append(commits, gitdomain.Commit{
					Message: gitdomain.CommitMessage(fmt.Sprintf("commit %d", i)),
					SHA:     gitdomain.SHA("1234567"),
				})
			}
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err = dialog.CommitsToBeam(commits, "target-branch", dialogTestInputs.Next())
			return err
		},
	}
}
