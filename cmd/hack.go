package cmd

import (
	"log"

	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
	"github.com/Originate/git-town/lib/steps"
	"github.com/Originate/git-town/lib/util"

	"github.com/spf13/cobra"
)

type HackFlags struct {
	Abort    bool
	Continue bool
}

var hackFlags HackFlags

var hackCmd = &cobra.Command{
	Use:   "hack",
	Short: "Create a new feature branch off the main development branch",
	Long:  `Create a new feature branch off the main development branch`,
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "hack",
			IsAbort:              hackFlags.Abort,
			IsContinue:           hackFlags.Continue,
			IsSkip:               false,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				targetBranchName := checkHackPreconditions(args)
				return getHackStepList(targetBranchName)
			},
		})
	},
}

func checkHackPreconditions(args []string) string {
	if len(args) == 0 {
		util.ExitWithErrorMessage("No branch name provided.")
	}
	targetBranchName := args[0]
	if config.HasRemote("origin") {
		fetchErr := script.RunCommand("git", "fetch", "--prune")
		if fetchErr != nil {
			log.Fatal(fetchErr)
		}
	}
	git.EnsureDoesNotHaveBranch(targetBranchName)
	return targetBranchName
}

func getHackStepList(targetBranchName string) steps.StepList {
	mainBranchName := config.GetMainBranch()
	stepList := steps.StepList{}
	stepList.AppendList(steps.GetSyncBranchSteps(mainBranchName))
	stepList.Append(steps.CreateAndCheckoutBranchStep{BranchName: targetBranchName, ParentBranchName: mainBranchName})
	if config.HasRemote("origin") && config.ShouldHackPush() {
		stepList.Append(steps.CreateTrackingBranchStep{BranchName: targetBranchName})
	}
	return steps.Wrap(stepList, steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
}

func init() {
	hackCmd.Flags().BoolVar(&hackFlags.Abort, "abort", false, "Abort a previous command that resulted in a conflict")
	hackCmd.Flags().BoolVar(&hackFlags.Continue, "continue", false, "Continue a previous command that resulted in a conflict")
	RootCmd.AddCommand(hackCmd)
}
