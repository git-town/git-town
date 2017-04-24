package cmd

import (
	"errors"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/script"
	"github.com/Originate/git-town/lib/steps"

	"github.com/spf13/cobra"
)

type prependConfig struct {
	InitialBranch string
	ParentBranch  string
	TargetBranch  string
}

var prependCommand = &cobra.Command{
	Use:   "prepend <branch>",
	Short: "Creates a new feature branch as the parent of the current branch",
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "prepend",
			IsAbort:              abortFlag,
			IsContinue:           continueFlag,
			IsSkip:               false,
			IsUndo:               undoFlag,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				config := checkPrependPreconditions(args)
				return getPrependStepList(config)
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

func checkPrependPreconditions(args []string) (result prependConfig) {
	result.InitialBranch = git.GetCurrentBranchName()
	result.TargetBranch = args[0]
	if git.HasRemote("origin") {
		script.Fetch()
	}
	git.EnsureDoesNotHaveBranch(result.TargetBranch)
	git.EnsureIsFeatureBranch(result.InitialBranch, "Only feature branches can have parent branches.")
	prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
	result.ParentBranch = git.GetParentBranch(result.InitialBranch)
	return
}

func getPrependStepList(config prependConfig) (result steps.StepList) {
	for _, branchName := range git.GetAncestorBranches(config.InitialBranch) {
		result.AppendList(steps.GetSyncBranchSteps(branchName))
	}
	result.Append(steps.CreateBranchStep{BranchName: config.TargetBranch, StartingPoint: config.ParentBranch})
	result.Append(steps.SetParentBranchStep{BranchName: config.TargetBranch, ParentBranchName: config.ParentBranch})
	result.Append(steps.SetParentBranchStep{BranchName: config.InitialBranch, ParentBranchName: config.TargetBranch})
	result.Append(steps.CheckoutBranchStep{BranchName: config.TargetBranch})
	if git.HasRemote("origin") && git.ShouldHackPush() {
		result.Append(steps.CreateTrackingBranchStep{BranchName: config.TargetBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return
}

func init() {
	prependCommand.Flags().BoolVar(&abortFlag, "abort", false, abortFlagDescription)
	prependCommand.Flags().BoolVar(&continueFlag, "continue", false, continueFlagDescription)
	prependCommand.Flags().BoolVar(&undoFlag, "undo", false, undoFlagDescription)
	RootCmd.AddCommand(prependCommand)
}
