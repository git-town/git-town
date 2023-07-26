package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/validate"
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
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                dryRun,
		Fetch:                 true,
		HandleUnfinishedState: true,
		OmitBranchNames:       false,
		ValidateIsOnline:      false,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	allBranches, initialBranch, err := execute.LoadBranches(&repo.Runner, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	config, err := determineSyncConfig(all, &repo.Runner, allBranches, initialBranch, repo.IsOffline)
	if err != nil {
		return err
	}
	stepList, err := syncBranchesSteps(config, &repo.Runner)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "sync",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &repo.Runner, nil, repo.RootDir)
}

type syncConfig struct {
	branchesToSync     git.BranchesSyncStatus
	hasOpenChanges     bool
	hasOrigin          bool
	hasUpstream        bool
	initialBranch      string
	isOffline          bool
	lineage            config.Lineage
	mainBranch         string
	previousBranch     string
	pullBranchStrategy config.PullBranchStrategy
	pushHook           bool
	shouldPushTags     bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

func determineSyncConfig(allFlag bool, run *git.ProdRunner, allBranchesSyncStatus git.BranchesSyncStatus, initialBranch string, isOffline bool) (*syncConfig, error) {
	previousBranch := run.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges, err := run.Backend.HasOpenChanges()
	if err != nil {
		return nil, err
	}
	hasOrigin, err := run.Backend.HasOrigin()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	var branchNamesToSync []string
	var shouldPushTags bool
	if allFlag {
		localBranches := allBranchesSyncStatus.LocalBranches().BranchNames()
		err = validate.KnowsBranchesAncestors(localBranches, mainBranch, &run.Backend)
		if err != nil {
			return nil, err
		}
		branchNamesToSync = localBranches
		shouldPushTags = true
	} else {
		err = validate.KnowsBranchAncestors(initialBranch, mainBranch, &run.Backend)
		if err != nil {
			return nil, err
		}
		branchNamesToSync = []string{initialBranch}
		shouldPushTags = !run.Config.IsFeatureBranch(initialBranch)
	}
	lineage := run.Config.Lineage()
	allBranchNamesToSync := lineage.BranchesAndAncestors(branchNamesToSync)
	syncStrategy, err := run.Config.SyncStrategy()
	if err != nil {
		return nil, err
	}
	pushHook, err := run.Config.PushHook()
	if err != nil {
		return nil, err
	}
	pullBranchStrategy, err := run.Config.PullBranchStrategy()
	if err != nil {
		return nil, err
	}
	shouldSyncUpstream, err := run.Config.ShouldSyncUpstream()
	if err != nil {
		return nil, err
	}
	hasUpstream, err := run.Backend.HasUpstream()
	if err != nil {
		return nil, err
	}
	branchesToSync, err := allBranchesSyncStatus.Select(allBranchNamesToSync)
	return &syncConfig{
		branchesToSync:     branchesToSync,
		hasOpenChanges:     hasOpenChanges,
		hasOrigin:          hasOrigin,
		hasUpstream:        hasUpstream,
		initialBranch:      initialBranch,
		isOffline:          isOffline,
		lineage:            lineage,
		mainBranch:         mainBranch,
		previousBranch:     previousBranch,
		pullBranchStrategy: pullBranchStrategy,
		pushHook:           pushHook,
		shouldPushTags:     shouldPushTags,
		shouldSyncUpstream: shouldSyncUpstream,
		syncStrategy:       syncStrategy,
	}, err
}

// syncBranchesSteps provides the step list for the "git sync" command.
func syncBranchesSteps(config *syncConfig, run *git.ProdRunner) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.branchesToSync {
		updateBranchSteps(&list, updateBranchStepsArgs{
			branch:             branch,
			config:             &run.Config,
			hasOrigin:          config.hasOrigin,
			hasUpstream:        config.hasUpstream,
			isOffline:          config.isOffline,
			lineage:            config.lineage,
			mainBranch:         config.mainBranch,
			pullBranchStrategy: config.pullBranchStrategy,
			pushBranch:         true,
			pushHook:           config.pushHook,
			shouldSyncUpstream: config.shouldSyncUpstream,
			syncStrategy:       config.syncStrategy,
		})
	}
	list.Add(&steps.CheckoutStep{Branch: config.initialBranch})
	if config.hasOrigin && config.shouldPushTags && !config.isOffline {
		list.Add(&steps.PushTagsStep{})
	}
	list.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranch,
		PreviousBranch:   config.previousBranch,
	})
	return list.Result()
}

// updateBranchSteps provides the steps to sync a particular branch.
func updateBranchSteps(list *runstate.StepListBuilder, args updateBranchStepsArgs) {
	isFeatureBranch := args.config.IsFeatureBranch(args.branch.Name)
	if !args.hasOrigin && !isFeatureBranch {
		return
	}
	list.Add(&steps.CheckoutStep{Branch: args.branch.Name})
	if isFeatureBranch {
		updateFeatureBranchSteps(list, args.branch, args.lineage, args.syncStrategy)
	} else {
		updatePerennialBranchSteps(list, updatePerennialBranchStepsArgs{
			branch:             args.branch,
			mainBranch:         args.mainBranch,
			pullBranchStrategy: args.pullBranchStrategy,
			shouldSyncUpstream: args.shouldSyncUpstream,
			hasUpstream:        args.hasUpstream,
		})
	}
	if args.pushBranch && args.hasOrigin && !args.isOffline {
		if !args.branch.HasTrackingBranch() {
			list.Add(&steps.CreateTrackingBranchStep{Branch: args.branch.Name})
			return
		}
		if !isFeatureBranch {
			list.Add(&steps.PushBranchStep{Branch: args.branch.Name})
			return
		}
		pushFeatureBranchSteps(list, args.branch.Name, args.syncStrategy, args.pushHook)
	}
}

type updateBranchStepsArgs struct {
	branch             git.BranchSyncStatus
	config             *git.RepoConfig
	hasOrigin          bool
	hasUpstream        bool
	isOffline          bool
	lineage            config.Lineage
	mainBranch         string
	pullBranchStrategy config.PullBranchStrategy
	pushBranch         bool
	pushHook           bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

func updateFeatureBranchSteps(list *runstate.StepListBuilder, branch git.BranchSyncStatus, lineage config.Lineage, syncStrategy config.SyncStrategy) {
	if branch.HasTrackingBranch() {
		syncBranchSteps(list, branch.TrackingBranch(), string(syncStrategy))
	}
	syncBranchSteps(list, lineage.Parent(branch.Name), string(syncStrategy))
}

func updatePerennialBranchSteps(list *runstate.StepListBuilder, args updatePerennialBranchStepsArgs) {
	if args.branch.HasTrackingBranch() {
		syncBranchSteps(list, args.branch.TrackingBranch(), string(args.pullBranchStrategy))
	}
	if args.mainBranch == args.branch.Name && args.hasUpstream && args.shouldSyncUpstream {
		list.Add(&steps.FetchUpstreamStep{Branch: args.mainBranch})
		list.Add(&steps.RebaseBranchStep{Branch: fmt.Sprintf("upstream/%s", args.mainBranch)})
	}
}

type updatePerennialBranchStepsArgs struct {
	branch             git.BranchSyncStatus
	mainBranch         string
	pullBranchStrategy config.PullBranchStrategy
	shouldSyncUpstream bool
	hasUpstream        bool
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
