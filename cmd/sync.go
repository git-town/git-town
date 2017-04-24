package cmd

import (
	"fmt"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/script"
	"github.com/Originate/git-town/lib/steps"

	"github.com/spf13/cobra"
)

type syncConfig struct {
	InitialBranch  string
	BranchesToSync []string
	ShouldPushTags bool
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Updates the current branch with all relevant changes",
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureIsRepository()
		steps.Run(steps.RunOptions{
			CanSkip: func() bool {
				return !(git.IsRebaseInProgress() && git.IsMainBranch(git.GetCurrentBranchName()))
			},
			Command:    "sync",
			IsAbort:    abortFlag,
			IsContinue: continueFlag,
			IsSkip:     skipFlag,
			IsUndo:     false,
			SkipMessageGenerator: func() string {
				return fmt.Sprintf("the sync of the '%s' branch", git.GetCurrentBranchName())
			},
			StepListGenerator: func() steps.StepList {
				return getSyncStepList(checkSyncPreconditions())
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
	},
}

func checkSyncPreconditions() (result syncConfig) {
	if git.HasRemote("origin") {
		script.Fetch()
	}
	result.InitialBranch = git.GetCurrentBranchName()
	if allFlag {
		branches := git.GetLocalBranchesWithMainBranchFirst()
		prompt.EnsureKnowsParentBranches(branches)
		result.BranchesToSync = branches
		result.ShouldPushTags = true
	} else if git.IsFeatureBranch(result.InitialBranch) {
		prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
		result.BranchesToSync = append(git.GetAncestorBranches(result.InitialBranch), result.InitialBranch)
	} else {
		result.BranchesToSync = []string{result.InitialBranch}
		result.ShouldPushTags = true
	}
	return
}

func getSyncStepList(config syncConfig) (result steps.StepList) {
	for _, branchName := range config.BranchesToSync {
		result.AppendList(steps.GetSyncBranchSteps(branchName))
	}
	result.Append(steps.CheckoutBranchStep{BranchName: config.InitialBranch})
	if git.HasRemote("origin") && config.ShouldPushTags {
		result.Append(steps.PushTagsStep{})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return
}

func init() {
	syncCmd.Flags().BoolVar(&allFlag, "all", false, "Sync all local branches")
	syncCmd.Flags().BoolVar(&abortFlag, "abort", false, abortFlagDescription)
	syncCmd.Flags().BoolVar(&continueFlag, "continue", false, continueFlagDescription)
	syncCmd.Flags().BoolVar(&skipFlag, "skip", false, "Continue a previous command by skipping the branch that resulted in a conflicted")
	RootCmd.AddCommand(syncCmd)
}
