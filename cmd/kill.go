package cmd

import (
	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/steps"
	"github.com/spf13/cobra"
)

type KillFlags struct {
	Undo bool
}

type KillConfig struct {
	InitialBranch       string
	IsTargetBranchLocal bool
	TargetBranch        string
}

var killFlags KillFlags

var killCommand = &cobra.Command{
	Use:   "kill [<branch>]",
	Short: "Removes an obsolete feature branch",
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "kill",
			IsAbort:              false,
			IsContinue:           false,
			IsSkip:               false,
			IsUndo:               killFlags.Undo,
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

	config.EnsureIsFeatureBranch(result.TargetBranch, "You can only kill feature branches.")

	result.IsTargetBranchLocal = git.HasLocalBranch(result.TargetBranch)
	if result.IsTargetBranchLocal {
		prompt.EnsureKnowsParentBranches([]string{result.TargetBranch})
	}

	if config.HasRemote("origin") {
		steps.FetchStep{}.Run()
	}

	if result.InitialBranch != result.TargetBranch {
		git.EnsureHasBranch(result.TargetBranch)
	}

	return
}

func getKillStepList(killConfig KillConfig) (result steps.StepList) {
	if killConfig.IsTargetBranchLocal {
		targetBranchParent := config.GetParentBranch(killConfig.TargetBranch)
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
		for _, child := range config.GetChildBranches(killConfig.TargetBranch) {
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
	killCommand.Flags().BoolVar(&killFlags.Undo, "undo", false, "Undo a previous command")
	RootCmd.AddCommand(killCommand)
}
