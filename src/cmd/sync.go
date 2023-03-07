package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/spf13/cobra"
)

func syncCmd(repo *git.ProdRepo) *cobra.Command {
	var allFlag bool
	var dryRunFlag bool
	syncCmd := cobra.Command{
		Use:   "sync",
		Short: "Updates the current branch with all relevant changes",
		Long: fmt.Sprintf(`Updates the current branch with all relevant changes

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
You can disable this by running "git config %s false".`, config.SyncUpstreamKey),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := determineSyncConfig(allFlag, repo)
			if err != nil {
				return err
			}
			stepList, err := syncBranchesSteps(config, repo)
			if err != nil {
				return err
			}
			runState := runstate.New("sync", stepList)
			return runstate.Execute(runState, repo, nil)
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			if dryRunFlag {
				currentBranch, err := repo.Silent.CurrentBranch()
				if err != nil {
					return err
				}
				repo.DryRun.Activate(currentBranch)
			}
			if err := validateIsConfigured(repo); err != nil {
				return err
			}
			exit, err := handleUnfinishedState(repo, nil)
			if err != nil {
				return err
			}
			if exit {
				os.Exit(0)
			}
			return nil
		},
		GroupID: "basic",
	}
	syncCmd.Flags().BoolVar(&allFlag, "all", false, "Sync all local branches")
	syncCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Print the commands but don't run them")
	return &syncCmd
}

type syncConfig struct {
	branchesToSync []string
	hasOrigin      bool
	initialBranch  string
	isOffline      bool
	shouldPushTags bool
}

func determineSyncConfig(allFlag bool, repo *git.ProdRepo) (*syncConfig, error) {
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return nil, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return nil, err
	}
	if hasOrigin && !isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return nil, err
		}
	}
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return nil, err
	}
	parentDialog := dialog.ParentBranches{}
	var branchesToSync []string
	var shouldPushTags bool
	if allFlag {
		branches, err := repo.Silent.LocalBranchesMainFirst()
		if err != nil {
			return nil, err
		}
		err = parentDialog.EnsureKnowsParentBranches(branches, repo)
		if err != nil {
			return nil, err
		}
		branchesToSync = branches
		shouldPushTags = true
	} else {
		err = parentDialog.EnsureKnowsParentBranches([]string{initialBranch}, repo)
		if err != nil {
			return nil, err
		}
		branchesToSync = append(repo.Config.AncestorBranches(initialBranch), initialBranch)
		shouldPushTags = !repo.Config.IsFeatureBranch(initialBranch)
	}
	return &syncConfig{
		branchesToSync: branchesToSync,
		hasOrigin:      hasOrigin,
		initialBranch:  initialBranch,
		isOffline:      isOffline,
		shouldPushTags: shouldPushTags,
	}, nil
}

// syncBranchesSteps provides the step list for the "git sync" command.
func syncBranchesSteps(config *syncConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.branchesToSync {
		updateBranchSteps(&list, branch, true, repo)
	}
	list.Add(&steps.CheckoutStep{Branch: config.initialBranch})
	if config.hasOrigin && config.shouldPushTags && !config.isOffline {
		list.Add(&steps.PushTagsStep{})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return list.Result()
}

// updateBranchSteps provides the steps to sync a particular branch.
func updateBranchSteps(list *runstate.StepListBuilder, branch string, pushBranch bool, repo *git.ProdRepo) {
	isFeatureBranch := repo.Config.IsFeatureBranch(branch)
	syncStrategy := list.SyncStrategy(repo.Config.SyncStrategy())
	hasOrigin := list.Bool(repo.Silent.HasOrigin())
	pushHook := list.Bool(repo.Config.PushHook())
	if !hasOrigin && !isFeatureBranch {
		return
	}
	list.Add(&steps.CheckoutStep{Branch: branch})
	if isFeatureBranch {
		updateFeatureBranchSteps(list, branch, repo)
	} else {
		updatePerennialBranchSteps(list, branch, repo)
	}
	isOffline := list.Bool(repo.Config.IsOffline())
	if pushBranch && hasOrigin && !isOffline {
		hasTrackingBranch := list.Bool(repo.Silent.HasTrackingBranch(branch))
		if !hasTrackingBranch {
			list.Add(&steps.CreateTrackingBranchStep{Branch: branch})
			return
		}
		if !isFeatureBranch {
			list.Add(&steps.PushBranchStep{Branch: branch})
			return
		}
		pushFeatureBranchSteps(list, branch, syncStrategy, pushHook)
	}
}

func updateFeatureBranchSteps(list *runstate.StepListBuilder, branch string, repo *git.ProdRepo) {
	syncStrategy := list.SyncStrategy(repo.Config.SyncStrategy())
	hasTrackingBranch := list.Bool(repo.Silent.HasTrackingBranch(branch))
	if hasTrackingBranch {
		syncBranchSteps(list, repo.Silent.TrackingBranch(branch), string(syncStrategy))
	}
	syncBranchSteps(list, repo.Config.ParentBranch(branch), string(syncStrategy))
}

func updatePerennialBranchSteps(list *runstate.StepListBuilder, branch string, repo *git.ProdRepo) {
	hasTrackingBranch := list.Bool(repo.Silent.HasTrackingBranch(branch))
	if hasTrackingBranch {
		pullBranchStrategy := list.PullBranchStrategy(repo.Config.PullBranchStrategy())
		syncBranchSteps(list, repo.Silent.TrackingBranch(branch), string(pullBranchStrategy))
	}
	mainBranch := repo.Config.MainBranch()
	hasUpstream := list.Bool(repo.Silent.HasRemote("upstream"))
	shouldSyncUpstream := list.Bool(repo.Config.ShouldSyncUpstream())
	if mainBranch == branch && hasUpstream && shouldSyncUpstream {
		list.Add(&steps.FetchUpstreamStep{Branch: mainBranch})
		list.Add(&steps.RebaseBranchStep{Branch: fmt.Sprintf("upstream/%s", mainBranch)})
	}
}

// syncBranchStep provides the steps to sync the given tracking branch into the current branch.
func syncBranchSteps(list *runstate.StepListBuilder, otherBranch, strategy string) {
	switch strategy {
	case "merge":
		list.Add(&steps.MergeStep{Branch: otherBranch})
	case "rebase":
		list.Add(&steps.RebaseBranchStep{Branch: otherBranch})
	default:
		list.Fail("unknown syncStrategy value: %q", strategy)
	}
}

func pushFeatureBranchSteps(list *runstate.StepListBuilder, branch string, syncStrategy config.SyncStrategy, pushHook bool) {
	switch syncStrategy {
	case config.SyncStrategyMerge:
		list.Add(&steps.PushBranchStep{Branch: branch, NoPushHook: !pushHook})
	case config.SyncStrategyRebase:
		list.Add(&steps.PushBranchStep{Branch: branch, ForceWithLease: true})
	default:
		list.Fail("unknown syncStrategy value: %q", syncStrategy)
	}
}
