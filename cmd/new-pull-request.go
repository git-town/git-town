package cmd

import (
	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/steps"
	"github.com/spf13/cobra"
)

type NewPullRequestConfig struct {
	InitialBranch  string
	BranchesToSync []string
}

type NewPullRequestFlags struct {
	Abort    bool
	Continue bool
}

var newPullRequestFlags NewPullRequestFlags

var newPullRequestCommand = &cobra.Command{
	Use:   "new-pull-request",
	Short: "Creates a new pull request",
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "new-pull-request",
			IsAbort:              newPullRequestFlags.Abort,
			IsContinue:           newPullRequestFlags.Continue,
			IsSkip:               false,
			IsUndo:               false,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				config := checkNewPullRequestPreconditions()
				return getNewPullRequestStepList(config)
			},
		})
	},
}

func checkNewPullRequestPreconditions() (result NewPullRequestConfig) {
	if config.HasRemote("origin") {
		steps.FetchStep{}.Run()
	}
	result.InitialBranch = git.GetCurrentBranchName()
	prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
	result.BranchesToSync = append(config.GetAncestorBranches(result.InitialBranch), result.InitialBranch)
	return
}

func getNewPullRequestStepList(config NewPullRequestConfig) steps.StepList {
	stepList := steps.StepList{}
	for _, branchName := range config.BranchesToSync {
		stepList.AppendList(steps.GetSyncBranchSteps(branchName))
	}
	stepList = steps.Wrap(stepList, steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	stepList.Append(steps.CreatePullRequestStep{BranchName: config.InitialBranch})
	return stepList
}

func init() {
	newPullRequestCommand.Flags().BoolVar(&newPullRequestFlags.Abort, "abort", false, "Abort a previous command that resulted in a conflict")
	newPullRequestCommand.Flags().BoolVar(&newPullRequestFlags.Continue, "continue", false, "Continue a previous command that resulted in a conflict")
	RootCmd.AddCommand(newPullRequestCommand)
}
