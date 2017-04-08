package cmd

import (
	"strings"

	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/steps"
	"github.com/Originate/git-town/lib/util"

	"github.com/spf13/cobra"
)

type ShipConfig struct {
	InitialBranch       string
	IsTargetBranchLocal bool
	TargetBranch        string
}

type ShipFlags struct {
	Abort         bool
	CommitMessage string
	Continue      bool
	Undo          bool
}

var shipFlags ShipFlags

var shipCmd = &cobra.Command{
	Use:   "ship",
	Short: "Deliver a completed feature branch",
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "ship",
			IsAbort:              shipFlags.Abort,
			IsContinue:           shipFlags.Continue,
			IsSkip:               false,
			IsUndo:               shipFlags.Undo,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				shipConfig := checkShipPreconditions(args)
				return getShipStepList(shipConfig)
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 1)
	},
}

func checkShipPreconditions(args []string) (result ShipConfig) {
	result.InitialBranch = git.GetCurrentBranchName()
	if len(args) == 0 {
		result.TargetBranch = result.InitialBranch
		git.EnsureDoesNotHaveOpenChanges("Did you mean to commit them before shipping?")
	} else {
		result.TargetBranch = args[0]
	}
	if config.HasRemote("origin") {
		steps.FetchStep{}.Run()
	}
	if result.TargetBranch != result.InitialBranch {
		git.EnsureHasBranch(result.TargetBranch)
	}
	config.EnsureIsFeatureBranch(result.TargetBranch, "Only feature branches can be shipped.")
	prompt.EnsureKnowsParentBranches([]string{result.TargetBranch})
	if config.GetParentBranch(result.TargetBranch) != config.GetMainBranch() {
		ancestors := config.GetAncestorBranches(result.TargetBranch)
		ancestorsWithoutMain := ancestors[1:]
		oldestAncestor := ancestorsWithoutMain[0]
		util.ExitWithErrorMessage(
			"Shipping this branch would ship "+strings.Join(ancestorsWithoutMain, ", ")+" as well.",
			"Please ship \""+oldestAncestor+"\" first.",
		)
	}
	return
}

func getShipStepList(shipConfig ShipConfig) steps.StepList {
	mainBranch := config.GetMainBranch()
	areInitialAndTargetDifferent := shipConfig.TargetBranch != shipConfig.InitialBranch
	stepList := steps.StepList{}
	stepList.AppendList(steps.GetSyncBranchSteps(mainBranch))
	stepList.Append(steps.CheckoutBranchStep{BranchName: shipConfig.TargetBranch})
	stepList.Append(steps.MergeTrackingBranchStep{})
	stepList.Append(steps.MergeBranchStep{BranchName: mainBranch})
	stepList.Append(steps.EnsureHasShippableChangesStep{BranchName: shipConfig.TargetBranch})
	stepList.Append(steps.CheckoutBranchStep{BranchName: mainBranch})
	stepList.Append(steps.SquashMergeBranchStep{BranchName: shipConfig.TargetBranch, CommitMessage: shipFlags.CommitMessage})
	if config.HasRemote("origin") {
		stepList.Append(steps.PushBranchStep{BranchName: mainBranch, Undoable: true})
	}
	childBranches := config.GetChildBranches(shipConfig.TargetBranch)
	if git.HasTrackingBranch(shipConfig.TargetBranch) && len(childBranches) == 0 {
		stepList.Append(steps.DeleteRemoteBranchStep{BranchName: shipConfig.TargetBranch})
	}
	stepList.Append(steps.DeleteLocalBranchStep{BranchName: shipConfig.TargetBranch})
	stepList.Append(steps.DeleteParentBranchStep{BranchName: shipConfig.TargetBranch})
	for _, child := range childBranches {
		stepList.Append(steps.SetParentBranchStep{BranchName: child, ParentBranchName: mainBranch})
	}
	stepList.Append(steps.DeleteAncestorBranchesStep{})
	if areInitialAndTargetDifferent {
		stepList.Append(steps.CheckoutBranchStep{BranchName: shipConfig.InitialBranch})
	}
	return steps.Wrap(stepList, steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: areInitialAndTargetDifferent})
}

func init() {
	shipCmd.Flags().BoolVar(&shipFlags.Abort, "abort", false, "Abort a previous command that resulted in a conflict")
	shipCmd.Flags().StringVarP(&shipFlags.CommitMessage, "message", "m", "", "Specify the commit message for the squash commit")
	shipCmd.Flags().BoolVar(&shipFlags.Continue, "continue", false, "Continue a previous command that resulted in a conflict")
	shipCmd.Flags().BoolVar(&shipFlags.Undo, "undo", false, "Undo a previous command")
	RootCmd.AddCommand(shipCmd)
}
