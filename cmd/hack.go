package cmd

import (
	"errors"
	"log"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/steps"

	"github.com/spf13/cobra"
)

var hackCmd = &cobra.Command{
	Use:   "hack <branch>",
	Short: "Create a new feature branch off the main development branch",
	Long:  `Create a new feature branch off the main development branch`,
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "hack",
			IsAbort:              abortFlag,
			IsContinue:           continueFlag,
			IsSkip:               false,
			IsUndo:               false,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				targetBranchName := checkHackPreconditions(args)
				return getHackStepList(targetBranchName)
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && !abortFlag && !continueFlag {
			return errors.New("no branch name provided")
		}
		return validateMaxArgs(args, 1)
	},
}

func checkHackPreconditions(args []string) string {
	targetBranchName := args[0]
	if git.HasRemote("origin") {
		err := steps.FetchStep{}.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	git.EnsureDoesNotHaveBranch(targetBranchName)
	return targetBranchName
}

func getHackStepList(targetBranchName string) steps.StepList {
	mainBranchName := git.GetMainBranch()
	stepList := steps.StepList{}
	stepList.AppendList(steps.GetSyncBranchSteps(mainBranchName))
	stepList.Append(steps.CreateAndCheckoutBranchStep{BranchName: targetBranchName, ParentBranchName: mainBranchName})
	if git.HasRemote("origin") && git.ShouldHackPush() {
		stepList.Append(steps.CreateTrackingBranchStep{BranchName: targetBranchName})
	}
	return steps.Wrap(stepList, steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
}

func init() {
	hackCmd.Flags().BoolVar(&abortFlag, "abort", false, "Abort a previous command that resulted in a conflict")
	hackCmd.Flags().BoolVar(&continueFlag, "continue", false, "Continue a previous command that resulted in a conflict")
	RootCmd.AddCommand(hackCmd)
}
