package cmd

import (
	"fmt"
	"os"

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
the main branch is synced with its upstream counterpart.
This can be disabled with "git config git-town.sync-upstream false".`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getSyncConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList := getSyncStepList(config)
		runState := steps.NewRunState("sync", stepList)
		err = steps.Run(runState)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			conditionallyActivateDryRun,
			validateIsConfigured,
			ensureIsNotInUnfinishedState,
		)
	},
}

func getSyncConfig() (result syncConfig, err error) {
	if git.HasRemote("origin") && !git.Config().IsOffline() {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}
	result.InitialBranch = git.GetCurrentBranchName()
	if allFlag {
		branches := git.GetLocalBranchesWithMainBranchFirst()
		prompt.EnsureKnowsParentBranches(branches)
		result.BranchesToSync = branches
		result.ShouldPushTags = true
	} else {
		prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
		result.BranchesToSync = append(git.Config().GetAncestorBranches(result.InitialBranch), result.InitialBranch)
		result.ShouldPushTags = !git.Config().IsFeatureBranch(result.InitialBranch)
	}
	return
}

func getSyncStepList(config syncConfig) (result steps.StepList) {
	for _, branchName := range config.BranchesToSync {
		result.AppendList(steps.GetSyncBranchSteps(branchName, true))
	}
	result.Append(&steps.CheckoutBranchStep{BranchName: config.InitialBranch})
	if git.HasRemote("origin") && config.ShouldPushTags && !git.Config().IsOffline() {
		result.Append(&steps.PushTagsStep{})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return
}

func init() {
	syncCmd.Flags().BoolVar(&allFlag, "all", false, "Sync all local branches")
	syncCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, dryRunFlagDescription)
	RootCmd.AddCommand(syncCmd)
}
