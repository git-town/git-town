package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
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
		Long:    long(syncDesc, fmt.Sprintf(syncHelp, config.KeySyncUpstream)),
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
	config, err := determineSyncConfig(all, &repo.Runner, repo.IsOffline)
	if err != nil {
		return err
	}
	stepList, err := syncBranchesSteps(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "sync",
		RunStepList: stepList,
	}
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &runState,
		Run:       &repo.Runner,
		Connector: nil,
		RootDir:   repo.RootDir,
	})
}

type syncConfig struct {
	branches           domain.Branches
	branchesToSync     domain.BranchInfos
	hasOpenChanges     bool
	remotes            config.Remotes
	isOffline          bool
	lineage            config.Lineage
	mainBranch         domain.LocalBranchName
	previousBranch     domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	pushHook           bool
	shouldPushTags     bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

func determineSyncConfig(allFlag bool, run *git.ProdRunner, isOffline bool) (*syncConfig, error) {
	branches, err := execute.LoadBranches(run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return nil, err
	}
	previousBranch := run.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges, err := run.Backend.HasOpenChanges()
	if err != nil {
		return nil, err
	}
	remotes, err := run.Backend.Remotes()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	var branchNamesToSync domain.LocalBranchNames
	var shouldPushTags bool
	lineage := run.Config.Lineage()
	var configUpdated bool
	if allFlag {
		localBranches := branches.All.LocalBranches()
		configUpdated, err = validate.KnowsBranchesAncestors(validate.KnowsBranchesAncestorsArgs{
			AllBranches: localBranches,
			Backend:     &run.Backend,
			BranchTypes: branches.Types,
			Lineage:     lineage,
			MainBranch:  mainBranch,
		})
		if err != nil {
			return nil, err
		}
		branchNamesToSync = localBranches.Names()
		shouldPushTags = true
	} else {
		configUpdated, err = validate.KnowsBranchAncestors(branches.Initial, validate.KnowsBranchAncestorsArgs{
			AllBranches:   branches.All,
			Backend:       &run.Backend,
			BranchTypes:   branches.Types,
			DefaultBranch: mainBranch,
			Lineage:       lineage,
			MainBranch:    mainBranch,
		})
		if err != nil {
			return nil, err
		}
	}
	if configUpdated {
		lineage = run.Config.Lineage() // reload after ancestry change
		branches.Types = run.Config.BranchTypes()
	}
	if !allFlag {
		branchNamesToSync = domain.LocalBranchNames{branches.Initial}
		if configUpdated {
			run.Config.Reload()
			branches.Types = run.Config.BranchTypes()
		}
		shouldPushTags = !branches.Types.IsFeatureBranch(branches.Initial)
	}
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
	branchesToSync, err := branches.All.Select(allBranchNamesToSync)
	return &syncConfig{
		branches:           branches,
		branchesToSync:     branchesToSync,
		hasOpenChanges:     hasOpenChanges,
		remotes:            remotes,
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
func syncBranchesSteps(config *syncConfig) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.branchesToSync {
		syncBranchSteps(&list, syncBranchStepsArgs{
			branch:             branch,
			branchTypes:        config.branches.Types,
			remotes:            config.remotes,
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
	list.Add(&steps.CheckoutStep{Branch: config.branches.Initial})
	if config.remotes.HasOrigin() && config.shouldPushTags && !config.isOffline {
		list.Add(&steps.PushTagsStep{})
	}
	list.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return list.Result()
}

// syncBranchSteps provides the steps to sync a particular branch.
func syncBranchSteps(list *runstate.StepListBuilder, args syncBranchStepsArgs) {
	isFeatureBranch := args.branchTypes.IsFeatureBranch(args.branch.Name)
	if !isFeatureBranch && !args.remotes.HasOrigin() {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	list.Add(&steps.CheckoutStep{Branch: args.branch.Name})
	if isFeatureBranch {
		syncFeatureBranchSteps(list, args.branch, args.lineage, args.syncStrategy)
	} else {
		syncPerennialBranchSteps(list, syncPerennialBranchStepsArgs{
			branch:             args.branch,
			mainBranch:         args.mainBranch,
			pullBranchStrategy: args.pullBranchStrategy,
			shouldSyncUpstream: args.shouldSyncUpstream,
			hasUpstream:        args.remotes.HasUpstream(),
		})
	}
	if args.pushBranch && args.remotes.HasOrigin() && !args.isOffline {
		switch {
		case !args.branch.HasTrackingBranch():
			list.Add(&steps.CreateTrackingBranchStep{Branch: args.branch.Name})
		case !isFeatureBranch:
			list.Add(&steps.PushBranchStep{Branch: args.branch.Name})
		default:
			pushFeatureBranchSteps(list, args.branch.Name, args.syncStrategy, args.pushHook)
		}
	}
}

type syncBranchStepsArgs struct {
	branch             domain.BranchInfo
	branchTypes        domain.BranchTypes
	remotes            config.Remotes
	isOffline          bool
	lineage            config.Lineage
	mainBranch         domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	pushBranch         bool
	pushHook           bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

func syncFeatureBranchSteps(list *runstate.StepListBuilder, branch domain.BranchInfo, lineage config.Lineage, syncStrategy config.SyncStrategy) {
	if branch.HasTrackingBranch() {
		updateCurrentFeatureBranchStep(list, branch.RemoteName.BranchName(), syncStrategy)
	}
	updateCurrentFeatureBranchStep(list, lineage.Parent(branch.Name).BranchName(), syncStrategy)
}

func syncPerennialBranchSteps(list *runstate.StepListBuilder, args syncPerennialBranchStepsArgs) {
	if args.branch.HasTrackingBranch() {
		updateCurrentPerennialBranchStep(list, args.branch.RemoteName, args.pullBranchStrategy)
	}
	if args.branch.Name == args.mainBranch && args.hasUpstream && args.shouldSyncUpstream {
		list.Add(&steps.FetchUpstreamStep{Branch: args.mainBranch})
		list.Add(&steps.RebaseBranchStep{Branch: domain.NewBranchName("upstream/" + args.mainBranch.String())})
	}
}

type syncPerennialBranchStepsArgs struct {
	branch             domain.BranchInfo
	mainBranch         domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	shouldSyncUpstream bool
	hasUpstream        bool
}

// updateCurrentFeatureBranchStep provides the step to update the current feature branch with changes from the given other branch.
func updateCurrentFeatureBranchStep(list *runstate.StepListBuilder, otherBranch domain.BranchName, strategy config.SyncStrategy) {
	switch strategy {
	case config.SyncStrategyMerge:
		list.Add(&steps.MergeStep{Branch: otherBranch})
	case config.SyncStrategyRebase:
		list.Add(&steps.RebaseBranchStep{Branch: otherBranch})
	default:
		list.Fail("unknown syncStrategy value: %q", strategy)
	}
}

// updateCurrentPerennialBranchStep provides the steps to update the current perennial branch with changes from the given other branch.
func updateCurrentPerennialBranchStep(list *runstate.StepListBuilder, otherBranch domain.RemoteBranchName, strategy config.PullBranchStrategy) {
	switch strategy {
	case config.PullBranchStrategyMerge:
		list.Add(&steps.MergeStep{Branch: otherBranch.BranchName()})
	case config.PullBranchStrategyRebase:
		list.Add(&steps.RebaseBranchStep{Branch: otherBranch.BranchName()})
	default:
		list.Fail("unknown syncStrategy value: %q", strategy)
	}
}

func pushFeatureBranchSteps(list *runstate.StepListBuilder, branch domain.LocalBranchName, syncStrategy config.SyncStrategy, pushHook bool) {
	switch syncStrategy {
	case config.SyncStrategyMerge:
		list.Add(&steps.PushBranchStep{Branch: branch, NoPushHook: !pushHook})
	case config.SyncStrategyRebase:
		list.Add(&steps.PushBranchStep{Branch: branch, ForceWithLease: true})
	default:
		list.Fail("unknown syncStrategy value: %q", syncStrategy)
	}
}
