package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
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
			entries := make(dialog.SwitchBranchEntries, amount)
			for a := range amount {
				entries[a] = dialog.SwitchBranchEntry{
					Branch:        gitdomain.NewLocalBranchName(fmt.Sprintf("branch-%d", a)),
					Indentation:   "",
					OtherWorktree: false,
					Type:          configdomain.BranchTypeFeatureBranch,
				}
			}
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, _, err = dialog.SwitchBranch(dialog.SwitchBranchArgs{
				Cursor:             0,
				DisplayBranchTypes: false,
				Entries:            entries,
				Inputs:             inputs,
				UncommittedChanges: false,
			})
			return err
		},
	}
}
