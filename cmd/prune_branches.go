package cmd

import (
	"log"

	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
	"github.com/Originate/git-town/lib/steps"
	"github.com/spf13/cobra"
)

type PruneBranchesFlags struct {
	Undo bool
}

var pruneBranchesFlags PruneBranchesFlags

var pruneBranchesCommand = &cobra.Command{
	Use:   "prune-branches",
	Short: "Deletes local branches whose tracking branch no longer exists",
	Long:  `Deletes local branches whose tracking branch no longer exists`,
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "prune-branches",
			IsAbort:              false,
			IsContinue:           false,
			IsSkip:               false,
			IsUndo:               pruneBranchesFlags.Undo,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				checkPruneBranchesPreconditions()
				return getPruneBranchesList()
			},
		})
	},
}

func checkPruneBranchesPreconditions() {
	if config.HasRemote("origin") {
		fetchErr := script.RunCommand("git", "fetch", "--prune")
		if fetchErr != nil {
			log.Fatal(fetchErr)
		}
	}
}

func getPruneBranchesList() (result steps.StepList) {
	initialBranchName := git.GetCurrentBranchName()
	for _, branchName := range git.GetLocalBranchesWithDeletedTrackingBranches() {
		if initialBranchName == branchName {
			result.Append(steps.CheckoutBranchStep{BranchName: config.GetMainBranch()})
		}

		parent := config.GetParentBranch(branchName)
		if parent != "" {
			for _, child := range config.GetChildBranches(branchName) {
				result.Append(steps.DeleteAncestorBranchesStep{BranchName: child})
				result.Append(steps.SetParentBranchStep{BranchName: child, ParentBranchName: parent})
			}
			result.Append(steps.DeleteParentBranchStep{BranchName: branchName})
		}

		result.Append(steps.DeleteLocalBranchStep{BranchName: branchName})
	}
	return
}

func init() {
	pruneBranchesCommand.Flags().BoolVar(&pruneBranchesFlags.Undo, "undo", false, "Undo a previous command")
	RootCmd.AddCommand(pruneBranchesCommand)
}
