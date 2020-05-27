package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/script"
	"github.com/git-town/git-town/src/steps"

	"github.com/spf13/cobra"
)

type syncConfig struct {
	initialBranch  string
	branchesToSync []string
	shouldPushTags bool
	hasOrigin      bool
	isOffline      bool
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

If the repository contains an "upstream" remote,
syncs the main branch with its upstream counterpart.
You can disable this by running "git config git-town.sync-upstream false".`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := git.NewProdRepo()
		config, err := getSyncConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList := getSyncStepList(config)
		runState := steps.NewRunState("sync", stepList)
		err = steps.Run(runState, repo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := git.ValidateIsRepository(); err != nil {
			return err
		}
		if err := conditionallyActivateDryRun(); err != nil {
			return err
		}
		if err := validateIsConfigured(); err != nil {
			return err
		}
		return ensureIsNotInUnfinishedState()
	},
}

func getSyncConfig() (result syncConfig, err error) {
	result.hasOrigin = git.HasRemote("origin")
	result.isOffline = git.Config().IsOffline()
	if result.hasOrigin && !result.isOffline {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}
	result.initialBranch = git.GetCurrentBranchName()
	if allFlag {
		branches := git.GetLocalBranchesWithMainBranchFirst()
		prompt.EnsureKnowsParentBranches(branches)
		result.branchesToSync = branches
		result.shouldPushTags = true
	} else {
		prompt.EnsureKnowsParentBranches([]string{result.initialBranch})
		result.branchesToSync = append(git.Config().GetAncestorBranches(result.initialBranch), result.initialBranch)
		result.shouldPushTags = !git.Config().IsFeatureBranch(result.initialBranch)
	}
	return
}

func getSyncStepList(config syncConfig) (result steps.StepList) {
	for _, branchName := range config.branchesToSync {
		result.AppendList(steps.GetSyncBranchSteps(branchName, true))
	}
	result.Append(&steps.CheckoutBranchStep{BranchName: config.initialBranch})
	if config.hasOrigin && config.shouldPushTags && !config.isOffline {
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
