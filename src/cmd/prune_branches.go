package cmd

import (
	"fmt"
	"os"

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
		err := checkPruneBranchesPreconditions()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList := getPruneBranchesStepList()
		runState := steps.NewRunState("prune-branches", stepList)
		err = steps.Run(runState)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
			git.Config().ValidateIsOnline,
		)
	},
}

func checkPruneBranchesPreconditions() error {
	if git.HasRemote("origin") {
		err := script.Fetch()
		if err != nil {
			return err
		}
	}
	return nil
}

func getPruneBranchesStepList() (result steps.StepList) {
	initialBranchName := git.GetCurrentBranchName()
	for _, branchName := range git.GetLocalBranchesWithDeletedTrackingBranches() {
		if initialBranchName == branchName {
			result.Append(&steps.CheckoutBranchStep{BranchName: git.Config().GetMainBranch()})
		}

		parent := git.Config().GetParentBranch(branchName)
		if parent != "" {
			for _, child := range git.Config().GetChildBranches(branchName) {
				result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: parent})
			}
			result.Append(&steps.DeleteParentBranchStep{BranchName: branchName})
		}

		if git.Config().IsPerennialBranch(branchName) {
			result.Append(&steps.RemoveFromPerennialBranches{BranchName: branchName})
		}

		result.Append(&steps.DeleteLocalBranchStep{BranchName: branchName})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: false, StashOpenChanges: false})
	return
}

func init() {
	RootCmd.AddCommand(pruneBranchesCommand)
}
