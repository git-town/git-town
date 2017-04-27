package cmd

import (
	"errors"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"

	"github.com/spf13/cobra"
)

var hackCmd = &cobra.Command{
	Use:   "hack <branch>",
	Short: "Creates a new feature branch off the main development branch",
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureIsRepository()
		prompt.EnsureIsConfigured()
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
		script.Fetch()
	}
	git.EnsureDoesNotHaveBranch(targetBranchName)
	return targetBranchName
}

func getHackStepList(targetBranchName string) (result steps.StepList) {
	mainBranchName := git.GetMainBranch()
	result.AppendList(steps.GetSyncBranchSteps(mainBranchName))
	result.Append(steps.CreateAndCheckoutBranchStep{BranchName: targetBranchName, ParentBranchName: mainBranchName})
	if git.HasRemote("origin") && git.ShouldHackPush() {
		result.Append(steps.CreateTrackingBranchStep{BranchName: targetBranchName})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return
}

func init() {
	hackCmd.Flags().BoolVar(&abortFlag, "abort", false, abortFlagDescription)
	hackCmd.Flags().BoolVar(&continueFlag, "continue", false, continueFlagDescription)
	RootCmd.AddCommand(hackCmd)
}
