package cmd

import (
	"log"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/gitconfig"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/steps"
	"github.com/spf13/cobra"
)

type killConfig struct {
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
			IsUndo:               UndoFlag,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				return getKillStepList(checkKillPreconditions(args))
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 1)
	},
}

func checkKillPreconditions(args []string) (result killConfig) {
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
		err := steps.FetchStep{}.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	if result.InitialBranch != result.TargetBranch {
		git.EnsureHasBranch(result.TargetBranch)
	}

	return
}

func getKillStepList(config killConfig) (result steps.StepList) {
	if config.IsTargetBranchLocal {
		targetBranchParent := gitconfig.GetParentBranch(config.TargetBranch)
		if git.HasTrackingBranch(config.TargetBranch) {
			result.Append(steps.DeleteRemoteBranchStep{BranchName: config.TargetBranch, IsTracking: true})
		}
		if config.InitialBranch == config.TargetBranch {
			if git.HasOpenChanges() {
				result.Append(steps.CommitOpenChangesStep{})
			}
			result.Append(steps.CheckoutBranchStep{BranchName: targetBranchParent})
		}
		result.Append(steps.DeleteLocalBranchStep{BranchName: config.TargetBranch, Force: true})
		for _, child := range gitconfig.GetChildBranches(config.TargetBranch) {
			result.Append(steps.SetParentBranchStep{BranchName: child, ParentBranchName: targetBranchParent})
		}
		result.Append(steps.DeleteParentBranchStep{BranchName: config.TargetBranch})
		result.Append(steps.DeleteAncestorBranchesStep{})
	} else {
		result.Append(steps.DeleteRemoteBranchStep{BranchName: config.TargetBranch, IsTracking: false})
	}
	return
}

func init() {
	killCommand.Flags().BoolVar(&UndoFlag, "undo", false, "Undo a previous command")
	RootCmd.AddCommand(killCommand)
}
