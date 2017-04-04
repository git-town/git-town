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
		git.EnsureDoesNotHaveUncommitedChanges("Did you mean to commit them before shipping?")
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
	stepList := steps.StepList{}
	stepList.AppendList(steps.GetSyncBranchSteps(mainBranch))
	stepList.Append(steps.CheckoutBranchStep{BranchName: shipConfig.TargetBranch})
	stepList.Append(steps.MergeTrackingBranchStep{})
	stepList.Append(steps.MergeBranchStep{BranchName: mainBranch})
	// TODO ensure_has_shippable_changes
	stepList.Append(steps.CheckoutBranchStep{BranchName: mainBranch})
	stepList.Append(steps.SquashMergeBranchStep{BranchName: shipConfig.TargetBranch})
	//stepList.Append(steps.CommitSquashMerge{CommitMessage: shipConfig.CommitMessage})
	// if [ "$HAS_REMOTE" = true ]; then
	//   echo "fetch"
	//   sync_branch_steps "$MAIN_BRANCH_NAME"
	// fi
	//
	// echo "checkout $target_branch_name"
	// echo "merge_tracking_branch"
	// echo "merge $MAIN_BRANCH_NAME"
	// echo "ensure_has_shippable_changes"
	// echo "checkout_main_branch"
	// echo "squash_merge $target_branch_name"
	// echo "commit_squash_merge $target_branch_name $commit_options"
	//
	// echo_if_true "push_branch $MAIN_BRANCH_NAME" "$HAS_REMOTE"
	//
	// if [ "$(has_tracking_branch "$target_branch_name")" = true ] &&
	//    [ "$(has_child_branches "$target_branch_name")" = false ]; then
	//   echo "delete_remote_branch $target_branch_name"
	// fi
	// echo "delete_local_branch $target_branch_name force"
	//
	// # update branch hierarchy information
	// echo "delete_parent_entry $target_branch_name"
	// echo_update_child_branches "$target_branch_name" "$MAIN_BRANCH_NAME"
	// echo "delete_all_ancestor_entries"
	//
	// if [ "$target_branch_name" != "$INITIAL_BRANCH_NAME" ]; then
	//   echo "checkout $INITIAL_BRANCH_NAME"
	// fi

	return steps.StepList{}
}

func init() {
	shipCmd.Flags().BoolVar(&shipFlags.Abort, "abort", false, "Abort a previous command that resulted in a conflict")
	shipCmd.Flags().StringVar(&shipFlags.CommitMessage, "m", "", "Specify the commit message for the squash commit")
	shipCmd.Flags().BoolVar(&shipFlags.Continue, "continue", false, "Continue a previous command that resulted in a conflict")
	shipCmd.Flags().BoolVar(&shipFlags.Undo, "undo", false, "Undo a previous command")
	RootCmd.AddCommand(shipCmd)
}
