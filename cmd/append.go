package cmd

import (
	"errors"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/script"
	"github.com/Originate/git-town/lib/steps"

	"github.com/spf13/cobra"
)

type appendConfig struct {
	InitialBranch string
	TargetBranch  string
}

var appendCommand = &cobra.Command{
	Use:   "append <branch>",
	Short: "Creates a new feature branch as a child of the current branch",
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "append",
			IsAbort:              abortFlag,
			IsContinue:           continueFlag,
			IsSkip:               false,
			IsUndo:               undoFlag,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				config := checkAppendPreconditions(args)
				return getAppendStepList(config)
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && !abortFlag && !continueFlag && !undoFlag {
			return errors.New("no branch name provided")
		}
		return validateMaxArgs(args, 1)
	},
}

func checkAppendPreconditions(args []string) (result appendConfig) {
	result.InitialBranch = git.GetCurrentBranchName()
	result.TargetBranch = args[0]
	if git.HasRemote("origin") {
		script.Fetch()
	}
	git.EnsureDoesNotHaveBranch(result.TargetBranch)
	prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
	return
}

func getAppendStepList(config appendConfig) steps.StepList {
	stepList := steps.StepList{}
	for _, branchName := range append(git.GetAncestorBranches(config.InitialBranch), config.InitialBranch) {
		stepList.AppendList(steps.GetSyncBranchSteps(branchName))
	}
	stepList.Append(steps.CreateBranchStep{BranchName: config.TargetBranch, StartingPoint: config.InitialBranch})
	stepList.Append(steps.SetParentBranchStep{BranchName: config.TargetBranch, ParentBranchName: config.InitialBranch})
	stepList.Append(steps.CheckoutBranchStep{BranchName: config.TargetBranch})
	if git.HasRemote("origin") && git.ShouldHackPush() {
		stepList.Append(steps.CreateTrackingBranchStep{BranchName: config.TargetBranch})
	}
	return steps.Wrap(stepList, steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
}

func init() {
	appendCommand.Flags().BoolVar(&abortFlag, "abort", false, "Abort a previous command that resulted in a conflict")
	appendCommand.Flags().BoolVar(&continueFlag, "continue", false, "Continue a previous command that resulted in a conflict")
	appendCommand.Flags().BoolVar(&undoFlag, "undo", false, "Undo a previous command")
	RootCmd.AddCommand(appendCommand)
}
