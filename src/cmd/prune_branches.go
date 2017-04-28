package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
	"github.com/spf13/cobra"
)

var pruneBranchesCommand = &cobra.Command{
	Use:   "prune-branches",
	Short: "Deletes local branches whose tracking branch no longer exists",
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureIsRepository()
		prompt.EnsureIsConfigured()
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "prune-branches",
			IsAbort:              false,
			IsContinue:           false,
			IsSkip:               false,
			IsUndo:               undoFlag,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				checkPruneBranchesPreconditions()
				return getPruneBranchesList()
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
	},
}

func checkPruneBranchesPreconditions() {
	if git.HasRemote("origin") {
		script.Fetch()
	}
}

func getPruneBranchesList() (result steps.StepList) {
	initialBranchName := git.GetCurrentBranchName()
	for _, branchName := range git.GetLocalBranchesWithDeletedTrackingBranches() {
		if initialBranchName == branchName {
			result.Append(steps.CheckoutBranchStep{BranchName: git.GetMainBranch()})
		}

		parent := git.GetParentBranch(branchName)
		if parent != "" {
			for _, child := range git.GetChildBranches(branchName) {
				result.Append(steps.SetParentBranchStep{BranchName: child, ParentBranchName: parent})
			}
			result.Append(steps.DeleteParentBranchStep{BranchName: branchName})
			result.Append(steps.DeleteAncestorBranchesStep{})
		}

		result.Append(steps.DeleteLocalBranchStep{BranchName: branchName})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: false, StashOpenChanges: false})
	return
}

func init() {
	pruneBranchesCommand.Flags().BoolVar(&undoFlag, "undo", false, undoFlagDescription)
	RootCmd.AddCommand(pruneBranchesCommand)
}
