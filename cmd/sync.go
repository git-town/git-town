package cmd

import (
	"fmt"

	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/steps"

	"github.com/spf13/cobra"
)

type SyncConfig struct {
	InitialBranch  string
	BranchesToSync []string
	ShouldPushTags bool
}

type SyncFlags struct {
	All      bool
	Abort    bool
	Continue bool
	Skip     bool
}

var syncFlags SyncFlags

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Updates the current branch with all relevant changes",
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip: func() bool {
				return !(git.IsRebaseInProgress() && config.IsMainBranch(git.GetCurrentBranchName()))
			},
			Command:    "sync",
			IsAbort:    syncFlags.Abort,
			IsContinue: syncFlags.Continue,
			IsSkip:     syncFlags.Skip,
			IsUndo:     false,
			SkipMessageGenerator: func() string {
				return fmt.Sprintf("the sync of the '%s' branch", git.GetCurrentBranchName())
			},
			StepListGenerator: func() steps.StepList {
				syncConfig := checkSyncPreconditions()
				return getSyncStepList(syncConfig)
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
	},
}

func checkSyncPreconditions() (result SyncConfig) {
	if config.HasRemote("origin") {
		steps.FetchStep{}.Run()
	}
	result.InitialBranch = git.GetCurrentBranchName()
	if syncFlags.All {
		branches := git.GetLocalBranchesWithMainBranchFirst()
		prompt.EnsureKnowsParentBranches(branches)
		result.BranchesToSync = branches
		result.ShouldPushTags = true
	} else if config.IsFeatureBranch(result.InitialBranch) {
		prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
		result.BranchesToSync = append(config.GetAncestorBranches(result.InitialBranch), result.InitialBranch)
	} else {
		result.BranchesToSync = []string{result.InitialBranch}
		result.ShouldPushTags = true
	}
	return
}

func getSyncStepList(syncConfig SyncConfig) steps.StepList {
	stepList := steps.StepList{}
	for _, branchName := range syncConfig.BranchesToSync {
		stepList.AppendList(steps.GetSyncBranchSteps(branchName))
	}
	stepList.Append(steps.CheckoutBranchStep{BranchName: syncConfig.InitialBranch})
	if config.HasRemote("origin") && syncConfig.ShouldPushTags {
		stepList.Append(steps.PushTagsStep{})
	}
	return steps.Wrap(stepList, steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
}

func init() {
	syncCmd.Flags().BoolVar(&syncFlags.All, "all", false, "Sync all local branches")
	syncCmd.Flags().BoolVar(&syncFlags.Abort, "abort", false, "Abort a previous command that resulted in a conflict")
	syncCmd.Flags().BoolVar(&syncFlags.Continue, "continue", false, "Continue a previous command that resulted in a conflict")
	syncCmd.Flags().BoolVar(&syncFlags.Skip, "skip", false, "Continue a previous command by skipping the branch that resulted in a conflicted")
	RootCmd.AddCommand(syncCmd)
}
