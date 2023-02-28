package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v7/src/cli"
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
You can disable this by running "git config %s false".`, config.SyncUpstream),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := determineSyncConfig(allFlag, repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := syncSteps(config, repo)
			if err != nil {
				cli.Exit(err)
			}
			runState := runstate.New("sync", stepList)
			err = runstate.Execute(runState, repo, nil)
			if err != nil {
				cli.Exit(err)
			}
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
	}
	syncCmd.Flags().BoolVar(&allFlag, "all", false, "Sync all local branches")
	syncCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Print the commands but don't run them")
	return &syncCmd
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

type syncConfig struct {
	branchesToSync []string
	hasOrigin      bool
	initialBranch  string
	isOffline      bool
	shouldPushTags bool
}

// syncSteps provides the step list for the "git sync" command.
func syncSteps(config *syncConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.branchesToSync {
		syncBranchSteps(&list, branch, true, repo)
	}
	list.Add(&steps.CheckoutBranchStep{Branch: config.initialBranch})
	if config.hasOrigin && config.shouldPushTags && !config.isOffline {
		list.Add(&steps.PushTagsStep{})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return list.Result()
}

// syncBranchSteps provides the steps to sync a particular branch.
func syncBranchSteps(list *runstate.StepListBuilder, branch string, pushBranch bool, repo *git.ProdRepo) {
	isFeatureBranch := repo.Config.IsFeatureBranch(branch)
	syncStrategy := repo.Config.SyncStrategy()
	hasOrigin := list.Bool(repo.Silent.HasOrigin())
	pushHook := list.Bool(repo.Config.PushHook())
	if !hasOrigin && !isFeatureBranch {
		return
	}
	list.Add(&steps.CheckoutBranchStep{Branch: branch})
	if isFeatureBranch {
		syncFeatureBranchSteps(list, branch, repo)
	} else {
		syncNonFeatureBranchSteps(list, branch, repo)
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

func syncFeatureBranchSteps(list *runstate.StepListBuilder, branch string, repo *git.ProdRepo) {
	syncStrategy := repo.Config.SyncStrategy()
	hasTrackingBranch := list.Bool(repo.Silent.HasTrackingBranch(branch))
	if hasTrackingBranch {
		syncTrackingBranchSteps(list, repo.Silent.TrackingBranch(branch), syncStrategy)
	}
	syncParentSteps(list, repo.Config.ParentBranch(branch), syncStrategy)
}

func syncNonFeatureBranchSteps(list *runstate.StepListBuilder, branch string, repo *git.ProdRepo) {
	hasTrackingBranch := list.Bool(repo.Silent.HasTrackingBranch(branch))
	if hasTrackingBranch {
		syncTrackingBranchSteps(list, repo.Silent.TrackingBranch(branch), repo.Config.PullBranchStrategy())
	}
	mainBranch := repo.Config.MainBranch()
	hasUpstream := list.Bool(repo.Silent.HasRemote("upstream"))
	shouldSyncUpstream := list.Bool(repo.Config.ShouldSyncUpstream())
	if mainBranch == branch && hasUpstream && shouldSyncUpstream {
		list.Add(&steps.FetchUpstreamStep{Branch: mainBranch})
		list.Add(&steps.RebaseBranchStep{Branch: fmt.Sprintf("upstream/%s", mainBranch)})
	}
}

// syncTrackingBranchStep provides the steps to sync the given tracking branch into the current branch.
func syncTrackingBranchSteps(list *runstate.StepListBuilder, trackingBranch, syncStrategy string) {
	switch syncStrategy {
	case "merge":
		list.Add(&steps.MergeBranchStep{Branch: trackingBranch})
	case "rebase":
		list.Add(&steps.RebaseBranchStep{Branch: trackingBranch})
	default:
		_ = list.Fail("unknown syncStrategy value: %q", syncStrategy)
	}
}

// syncParentSteps provides the steps to sync the given parent branch into the current branch.
func syncParentSteps(list *runstate.StepListBuilder, parentBranch, syncStrategy string) {
	switch syncStrategy {
	case "merge":
		list.Add(&steps.MergeBranchStep{Branch: parentBranch})
	case "rebase":
		list.Add(&steps.RebaseBranchStep{Branch: parentBranch})
	default:
		_ = list.Fail("unknown syncStrategy value: %q", syncStrategy)
	}
}

func pushFeatureBranchSteps(list *runstate.StepListBuilder, branch, syncStrategy string, pushHook bool) {
	switch syncStrategy {
	case "merge":
		list.Add(&steps.PushBranchStep{Branch: branch, NoPushHook: !pushHook})
	case "rebase":
		list.Add(&steps.PushBranchStep{Branch: branch, ForceWithLease: true})
	default:
		_ = list.Fail("unknown syncStrategy value: %q", syncStrategy)
	}
}
