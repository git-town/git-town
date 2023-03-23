package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/flags"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

const syncDesc = "Updates the current branch with all relevant changes"

const syncHelp = `
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
You can disable this by running "git config %s false".`

func syncCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addAllFlag, readAllFlag := flags.Bool("all", "a", "Sync all local branches")
	cmd := cobra.Command{
		Use:     "sync",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   syncDesc,
		Long:    long(syncDesc, fmt.Sprintf(syncHelp, config.SyncUpstreamKey)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return sync(readAllFlag(cmd), readDryRunFlag(cmd), readDebugFlag(cmd))
		},
	}
	addAllFlag(&cmd)
	addDebugFlag(&cmd)
	addDryRunFlag(&cmd)
	return &cmd
}

func sync(all, dryRun, debug bool) error {
	run, exit, err := LoadProdRunner(RepoArgs{
		debug:                 debug,
		dryRun:                dryRun,
		handleUnfinishedState: true,
		validateGitversion:    true,
		validateIsRepository:  true,
		validateIsConfigured:  true,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineSyncConfig(all, &run)
	if err != nil {
		return err
	}
	stepList, err := syncBranchesSteps(config, &run)
	if err != nil {
		return err
	}
	runState := runstate.New("sync", stepList)
	return runstate.Execute(runState, &run, nil)
}

type syncConfig struct {
	branchesToSync []string
	hasOrigin      bool
	initialBranch  string
	isOffline      bool
	mainBranch     string
	shouldPushTags bool
}

func determineSyncConfig(allFlag bool, run *git.ProdRunner) (*syncConfig, error) {
	hasOrigin, err := run.Backend.HasOrigin()
	if err != nil {
		return nil, err
	}
	isOffline, err := run.Config.IsOffline()
	if err != nil {
		return nil, err
	}
	if hasOrigin && !isOffline {
		err := run.Frontend.Fetch()
		if err != nil {
			return nil, err
		}
	}
	initialBranch, err := run.Backend.CurrentBranch()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	var branchesToSync []string
	var shouldPushTags bool
	if allFlag {
		branches, err := run.Backend.LocalBranchesMainFirst(mainBranch)
		if err != nil {
			return nil, err
		}
		err = validate.KnowsBranchesAncestry(branches, &run.Backend)
		if err != nil {
			return nil, err
		}
		branchesToSync = branches
		shouldPushTags = true
	} else {
		err = validate.KnowsBranchAncestry(initialBranch, run.Config.MainBranch(), &run.Backend)
		if err != nil {
			return nil, err
		}
		branchesToSync = append(run.Config.AncestorBranches(initialBranch), initialBranch)
		shouldPushTags = !run.Config.IsFeatureBranch(initialBranch)
	}
	return &syncConfig{
		branchesToSync: branchesToSync,
		hasOrigin:      hasOrigin,
		initialBranch:  initialBranch,
		isOffline:      isOffline,
		mainBranch:     mainBranch,
		shouldPushTags: shouldPushTags,
	}, nil
}

// syncBranchesSteps provides the step list for the "git sync" command.
func syncBranchesSteps(config *syncConfig, run *git.ProdRunner) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.branchesToSync {
		updateBranchSteps(&list, branch, true, run)
	}
	list.Add(&steps.CheckoutStep{Branch: config.initialBranch})
	if config.hasOrigin && config.shouldPushTags && !config.isOffline {
		list.Add(&steps.PushTagsStep{})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, &run.Backend, config.mainBranch)
	return list.Result()
}

// updateBranchSteps provides the steps to sync a particular branch.
func updateBranchSteps(list *runstate.StepListBuilder, branch string, pushBranch bool, run *git.ProdRunner) {
	isFeatureBranch := run.Config.IsFeatureBranch(branch)
	syncStrategy := list.SyncStrategy(run.Config.SyncStrategy())
	hasOrigin := list.Bool(run.Backend.HasOrigin())
	pushHook := list.Bool(run.Config.PushHook())
	if !hasOrigin && !isFeatureBranch {
		return
	}
	list.Add(&steps.CheckoutStep{Branch: branch})
	if isFeatureBranch {
		updateFeatureBranchSteps(list, branch, run)
	} else {
		updatePerennialBranchSteps(list, branch, run)
	}
	isOffline := list.Bool(run.Config.IsOffline())
	if pushBranch && hasOrigin && !isOffline {
		hasTrackingBranch := list.Bool(run.Backend.HasTrackingBranch(branch))
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

func updateFeatureBranchSteps(list *runstate.StepListBuilder, branch string, run *git.ProdRunner) {
	syncStrategy := list.SyncStrategy(run.Config.SyncStrategy())
	hasTrackingBranch := list.Bool(run.Backend.HasTrackingBranch(branch))
	if hasTrackingBranch {
		syncBranchSteps(list, run.Backend.TrackingBranch(branch), string(syncStrategy))
	}
	syncBranchSteps(list, run.Config.ParentBranch(branch), string(syncStrategy))
}

func updatePerennialBranchSteps(list *runstate.StepListBuilder, branch string, run *git.ProdRunner) {
	hasTrackingBranch := list.Bool(run.Backend.HasTrackingBranch(branch))
	if hasTrackingBranch {
		pullBranchStrategy := list.PullBranchStrategy(run.Config.PullBranchStrategy())
		syncBranchSteps(list, run.Backend.TrackingBranch(branch), string(pullBranchStrategy))
	}
	mainBranch := run.Config.MainBranch()
	hasUpstream := list.Bool(run.Backend.HasRemote("upstream"))
	shouldSyncUpstream := list.Bool(run.Config.ShouldSyncUpstream())
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
