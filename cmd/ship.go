package cmd

import (
	"strings"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/script"
	"github.com/Originate/git-town/lib/steps"
	"github.com/Originate/git-town/lib/util"

	"github.com/spf13/cobra"
)

type shipConfig struct {
	InitialBranch       string
	IsTargetBranchLocal bool
	TargetBranch        string
}

var commitMessage string

var shipCmd = &cobra.Command{
	Use:   "ship",
	Short: "Deliver a completed feature branch",
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "ship",
			IsAbort:              abortFlag,
			IsContinue:           continueFlag,
			IsSkip:               false,
			IsUndo:               undoFlag,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				config := checkShipPreconditions(args)
				return getShipStepList(config)
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 1)
	},
}

func checkShipPreconditions(args []string) (result shipConfig) {
	result.InitialBranch = git.GetCurrentBranchName()
	if len(args) == 0 {
		result.TargetBranch = result.InitialBranch
		git.EnsureDoesNotHaveOpenChanges("Did you mean to commit them before shipping?")
	} else {
		result.TargetBranch = args[0]
	}
	if git.HasRemote("origin") {
		script.Fetch()
	}
	if result.TargetBranch != result.InitialBranch {
		git.EnsureHasBranch(result.TargetBranch)
	}
	git.EnsureIsFeatureBranch(result.TargetBranch, "Only feature branches can be shipped.")
	prompt.EnsureKnowsParentBranches([]string{result.TargetBranch})
	ensureParentBranchIsMainBranch(result.TargetBranch)
	return
}

func ensureParentBranchIsMainBranch(branchName string) {
	if git.GetParentBranch(branchName) != git.GetMainBranch() {
		ancestors := git.GetAncestorBranches(branchName)
		ancestorsWithoutMain := ancestors[1:]
		oldestAncestor := ancestorsWithoutMain[0]
		util.ExitWithErrorMessage(
			"Shipping this branch would ship "+strings.Join(ancestorsWithoutMain, ", ")+" as well.",
			"Please ship \""+oldestAncestor+"\" first.",
		)
	}
}

func getShipStepList(config shipConfig) (result steps.StepList) {
	mainBranch := git.GetMainBranch()
	areInitialAndTargetDifferent := config.TargetBranch != config.InitialBranch
	result.AppendList(steps.GetSyncBranchSteps(mainBranch))
	result.Append(steps.CheckoutBranchStep{BranchName: config.TargetBranch})
	result.Append(steps.MergeTrackingBranchStep{})
	result.Append(steps.MergeBranchStep{BranchName: mainBranch})
	result.Append(steps.EnsureHasShippableChangesStep{BranchName: config.TargetBranch})
	result.Append(steps.CheckoutBranchStep{BranchName: mainBranch})
	result.Append(steps.SquashMergeBranchStep{BranchName: config.TargetBranch, CommitMessage: commitMessage})
	if git.HasRemote("origin") {
		result.Append(steps.PushBranchStep{BranchName: mainBranch, Undoable: true})
	}
	childBranches := git.GetChildBranches(config.TargetBranch)
	if git.HasTrackingBranch(config.TargetBranch) && len(childBranches) == 0 {
		result.Append(steps.DeleteRemoteBranchStep{BranchName: config.TargetBranch, IsTracking: true})
	}
	result.Append(steps.DeleteLocalBranchStep{BranchName: config.TargetBranch})
	result.Append(steps.DeleteParentBranchStep{BranchName: config.TargetBranch})
	for _, child := range childBranches {
		result.Append(steps.SetParentBranchStep{BranchName: child, ParentBranchName: mainBranch})
	}
	result.Append(steps.DeleteAncestorBranchesStep{})
	if areInitialAndTargetDifferent {
		result.Append(steps.CheckoutBranchStep{BranchName: config.InitialBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: areInitialAndTargetDifferent})
	return
}

func init() {
	shipCmd.Flags().BoolVar(&abortFlag, "abort", false, "Abort a previous command that resulted in a conflict")
	shipCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Specify the commit message for the squash commit")
	shipCmd.Flags().BoolVar(&continueFlag, "continue", false, "Continue a previous command that resulted in a conflict")
	shipCmd.Flags().BoolVar(&undoFlag, "undo", false, "Undo a previous command")
	RootCmd.AddCommand(shipCmd)
}
