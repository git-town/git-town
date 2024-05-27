package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
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
			localBranches := gitdomain.LocalBranchNames{}
			branchInfos := gitdomain.BranchInfos{}
			for i := 0; i < int(amount); i++ {
				branchName := gitdomain.NewLocalBranchName(fmt.Sprintf("branch-%d", i))
				localBranches = append(localBranches, branchName)
				branchInfos = append(branchInfos, gitdomain.BranchInfo{LocalName: Some(branchName), SyncStatus: gitdomain.SyncStatusLocalOnly}) //exhaustruct:ignore
			}
			lineage := configdomain.Lineage{}
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err = dialog.SwitchBranch(localBranches, gitdomain.NewLocalBranchName("branch-2"), lineage, branchInfos, true, dialogTestInputs.Next())
			return err
		},
	}
}
