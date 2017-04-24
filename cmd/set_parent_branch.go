package cmd

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/spf13/cobra"
)

var setParentBranchCommand = &cobra.Command{
	Use:   "set-parent-branch <child_branch> <parent_branch>",
	Short: "Updates a branch's parent",
	Run: func(cmd *cobra.Command, args []string) {
		setParentBranch(args[0], args[1])
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgsCount(args, 2)
	},
}

func setParentBranch(childBranch, parentBranch string) {
	git.EnsureIsRepository()
	git.EnsureHasBranch(childBranch)
	git.EnsureHasBranch(parentBranch)
	git.SetParentBranch(childBranch, parentBranch)
	git.DeleteAllAncestorBranches()
}

func init() {
	RootCmd.AddCommand(setParentBranchCommand)
}
