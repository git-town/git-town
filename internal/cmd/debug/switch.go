package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
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
			branchInfos := gitdomain.BranchInfos{}
			for i := range amount {
				branchName := gitdomain.NewLocalBranchName(fmt.Sprintf("branch-%d", i))
				branchInfos = append(branchInfos, gitdomain.BranchInfo{LocalName: Some(branchName), SyncStatus: gitdomain.SyncStatusLocalOnly}) //exhaustruct:ignore
			}
			lineage := configdomain.Lineage{}
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			branchTypes := []configdomain.BranchType{}
			branchesAndTypes := configdomain.BranchesAndTypes{}
			_, _, err = dialog.SwitchBranch(branchTypes, branchesAndTypes, gitdomain.NewLocalBranchName("branch-2"), lineage, branchInfos, configdomain.DefaultBranchType{BranchType: configdomain.BranchTypeFeatureBranch}, true, dialogTestInputs.Next())
			return err
		},
	}
}
