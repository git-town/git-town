package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

func switchBranch() *cobra.Command {
	return &cobra.Command{
		Use:  "switch <number of branches>",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			amount, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			entries := make([]dialog.SwitchBranchEntry, amount)
			for a := range amount {
				entries[a] = dialog.SwitchBranchEntry{
					Branch:        gitdomain.NewLocalBranchName(fmt.Sprintf("branch-%d", i)),
					Indentation:   "",
					OtherWorktree: false,
				}
			}
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err = dialog.SwitchBranch(entries, 0, false, dialogTestInputs.Next())
			return err
		},
	}
}
