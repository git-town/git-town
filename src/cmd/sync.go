package cmd

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"

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
	Long: `Updates the current branch with all relevant changes

Synchronizes the current branch with the rest of the world.

When run on a feature branch
- syncs all ancestor branches
- pulls updates for the current branch
- merges the parent branch into the current branch
- pushes the current branch

When run on the main branch or a perennial branch
- pulls and pushes updates for the current branch
- pushes tags

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.`,
	Run: func(cmd *cobra.Command, args []string) {
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
				return getSyncStepList(getSyncConfig())
			},
		})
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			conditionallyActivateDryRun,
			validateIsConfigured,
		)
	},
}

func getSyncConfig() (result syncConfig) {
	if git.HasRemote("origin") && !git.IsOffline() {
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
		result.AppendList(steps.GetSyncBranchSteps(branchName, true))
	}
	result.Append(&steps.CheckoutBranchStep{BranchName: config.InitialBranch})
	if git.HasRemote("origin") && config.ShouldPushTags && !git.IsOffline() {
		result.Append(&steps.PushTagsStep{})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return
}

func init() {
	syncCmd.Flags().BoolVar(&allFlag, "all", false, "Sync all local branches")
	syncCmd.Flags().BoolVar(&abortFlag, "abort", false, abortFlagDescription)
	syncCmd.Flags().BoolVar(&continueFlag, "continue", false, continueFlagDescription)
	syncCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, dryRunFlagDescription)
	syncCmd.Flags().BoolVar(&skipFlag, "skip", false, "Continue a previous command by skipping the branch that resulted in a conflicted")
	RootCmd.AddCommand(syncCmd)
}
