package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"
	"github.com/spf13/cobra"
)

var pruneBranchesCommand = &cobra.Command{
	Use:   "prune-branches",
	Short: "Deletes local branches whose tracking branch no longer exists",
	Long: `Deletes local branches whose tracking branch no longer exists

Deletes branches whose tracking branch no longer exists from the local repository.
This usually means the branch was shipped or killed on another machine.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkPruneBranchesPreconditions()
		stepList := getPruneBranchesStepList()
		runState := steps.NewRunState("prune-branches", stepList)
		steps.Run(runState)
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
			git.ValidateIsOnline,
		)
	},
}

func checkPruneBranchesPreconditions() {
	if git.HasRemote("origin") {
		script.Fetch()
	}
}

func getPruneBranchesStepList() (result steps.StepList) {
	initialBranchName := git.GetCurrentBranchName()
	for _, branchName := range git.GetLocalBranchesWithDeletedTrackingBranches() {
		if initialBranchName == branchName {
			result.Append(&steps.CheckoutBranchStep{BranchName: git.GetMainBranch()})
		}

		parent := git.GetParentBranch(branchName)
		if parent != "" {
			for _, child := range git.GetChildBranches(branchName) {
				result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: parent})
			}
			result.Append(&steps.DeleteParentBranchStep{BranchName: branchName})
		}

		result.Append(&steps.DeleteLocalBranchStep{BranchName: branchName})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: false, StashOpenChanges: false})
	return
}

func init() {
	RootCmd.AddCommand(pruneBranchesCommand)
}
