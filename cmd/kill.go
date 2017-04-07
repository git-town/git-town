package cmd

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/gitconfig"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/steps"
	"github.com/spf13/cobra"
)

type KillConfig struct {
	InitialBranch       string
	IsTargetBranchLocal bool
	TargetBranch        string
}

var killCommand = &cobra.Command{
	Use:   "kill [<branch>]",
	Short: "Removes an obsolete feature branch",
	Long:  "Removes an obsolete feature branch",
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "kill",
			IsAbort:              false,
			IsContinue:           false,
			IsSkip:               false,
			IsUndo:               undoFlag,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				killConfig := checkKillPreconditions(args)
				return getKillStepList(killConfig)
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 1)
	},
}

func checkKillPreconditions(args []string) (result KillConfig) {
	result.InitialBranch = git.GetCurrentBranchName()

	if len(args) == 0 {
		result.TargetBranch = result.InitialBranch
	} else {
		result.TargetBranch = args[0]
	}

	gitconfig.EnsureIsFeatureBranch(result.TargetBranch, "You can only kill feature branches.")

	result.IsTargetBranchLocal = git.HasLocalBranch(result.TargetBranch)
	if result.IsTargetBranchLocal {
		prompt.EnsureKnowsParentBranches([]string{result.TargetBranch})
	}

	if gitconfig.HasRemote("origin") {
		steps.FetchStep{}.Run()
	}

	if result.InitialBranch != result.TargetBranch {
		git.EnsureHasBranch(result.TargetBranch)
	}

	return
}

func getKillStepList(killConfig KillConfig) (result steps.StepList) {
	if killConfig.IsTargetBranchLocal {
		targetBranchParent := gitconfig.GetParentBranch(killConfig.TargetBranch)
		if git.HasTrackingBranch(killConfig.TargetBranch) {
			result.Append(steps.DeleteRemoteBranchStep{BranchName: killConfig.TargetBranch, IsTracking: true})
		}
		if killConfig.InitialBranch == killConfig.TargetBranch {
			if git.HasOpenChanges() {
				result.Append(steps.CommitOpenChangesStep{})
			}
			result.Append(steps.CheckoutBranchStep{BranchName: targetBranchParent})
		}
		result.Append(steps.DeleteLocalBranchStep{BranchName: killConfig.TargetBranch, Force: true})
		for _, child := range gitconfig.GetChildBranches(killConfig.TargetBranch) {
			result.Append(steps.SetParentBranchStep{BranchName: child, ParentBranchName: targetBranchParent})
		}
		result.Append(steps.DeleteParentBranchStep{BranchName: killConfig.TargetBranch})
		result.Append(steps.DeleteAncestorBranchesStep{})
	} else {
		result.Append(steps.DeleteRemoteBranchStep{BranchName: killConfig.TargetBranch, IsTracking: false})
	}
	return
}

func init() {
	killCommand.Flags().BoolVar(&undoFlag, "undo", false, "Undo a previous command")
	RootCmd.AddCommand(killCommand)
}
