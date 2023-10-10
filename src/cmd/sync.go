package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/runvm"
	"github.com/git-town/git-town/v9/src/step"
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
			return executeSync(readAllFlag(cmd), readDryRunFlag(cmd), readDebugFlag(cmd))
		},
	}
	addAllFlag(&cmd)
	addDebugFlag(&cmd)
	addDryRunFlag(&cmd)
	return &cmd
}

func executeSync(all, dryRun, debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           dryRun,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineSyncConfig(all, repo, debug)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "sync",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunSteps: syncBranchesSteps(config.branchesToSync, syncBranchStepsArgs{
			branches:           config.branches,
			branchTypes:        config.branches.Types,
			remotes:            config.remotes,
			hasOpenChanges:     config.hasOpenChanges,
			isOffline:          config.isOffline,
			lineage:            config.lineage,
			mainBranch:         config.mainBranch,
			previousBranch:     config.previousBranch,
			pullBranchStrategy: config.pullBranchStrategy,
			pushBranch:         config.pushBranch,
			pushHook:           config.pushHook,
			shouldPushTags:     config.shouldPushTags,
			shouldSyncUpstream: config.shouldSyncUpstream,
			syncStrategy:       config.syncStrategy,
		}),
	}
	return runvm.Execute(runvm.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               nil,
		Debug:                   debug,
		Lineage:                 config.lineage,
		NoPushHook:              !config.pushHook,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
	})
}

type syncConfig struct {
	branches           domain.Branches
	branchesToSync     domain.BranchInfos
	hasOpenChanges     bool
	remotes            domain.Remotes
	isOffline          bool
	lineage            config.Lineage
	mainBranch         domain.LocalBranchName
	previousBranch     domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	pushBranch         bool // TODO: is this used?
	pushHook           bool
	shouldPushTags     bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

func determineSyncConfig(allFlag bool, repo *execute.OpenRepoResult, debug bool) (*syncConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage()
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Debug:                 debug,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSnapshot, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	remotes, err := repo.Runner.Backend.Remotes()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	var branchNamesToSync domain.LocalBranchNames
	var shouldPushTags bool
	var configUpdated bool
	if allFlag {
		localBranches := branches.All.LocalBranches()
		configUpdated, err = validate.KnowsBranchesAncestors(validate.KnowsBranchesAncestorsArgs{
			AllBranches: localBranches,
			Backend:     &repo.Runner.Backend,
			BranchTypes: branches.Types,
			MainBranch:  mainBranch,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSnapshot, false, err
		}
		branchNamesToSync = localBranches.Names()
		shouldPushTags = true
	} else {
		configUpdated, err = validate.KnowsBranchAncestors(branches.Initial, validate.KnowsBranchAncestorsArgs{
			AllBranches:   branches.All,
			Backend:       &repo.Runner.Backend,
			BranchTypes:   branches.Types,
			DefaultBranch: mainBranch,
			MainBranch:    mainBranch,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSnapshot, false, err
		}
	}
	if configUpdated {
		lineage = repo.Runner.Config.Lineage() // reload after ancestry change
		branches.Types = repo.Runner.Config.BranchTypes()
	}
	if !allFlag {
		branchNamesToSync = domain.LocalBranchNames{branches.Initial}
		if configUpdated {
			repo.Runner.Config.Reload()
			branches.Types = repo.Runner.Config.BranchTypes()
		}
		shouldPushTags = !branches.Types.IsFeatureBranch(branches.Initial)
	}
	allBranchNamesToSync := lineage.BranchesAndAncestors(branchNamesToSync)
	syncStrategy, err := repo.Runner.Config.SyncStrategy()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	pullBranchStrategy, err := repo.Runner.Config.PullBranchStrategy()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	shouldSyncUpstream, err := repo.Runner.Config.ShouldSyncUpstream()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	branchesToSync, err := branches.All.Select(allBranchNamesToSync)
	return &syncConfig{
		branches:           branches,
		branchesToSync:     branchesToSync,
		hasOpenChanges:     repoStatus.OpenChanges,
		remotes:            remotes,
		isOffline:          repo.IsOffline,
		lineage:            lineage,
		mainBranch:         mainBranch,
		previousBranch:     previousBranch,
		pullBranchStrategy: pullBranchStrategy,
		pushHook:           pushHook,
		shouldPushTags:     shouldPushTags,
		shouldSyncUpstream: shouldSyncUpstream,
		syncStrategy:       syncStrategy,
	}, branchesSnapshot, stashSnapshot, false, err
}

// syncBranchesSteps provides the step list for the "git sync" command.
func syncBranchesSteps(branchesToSync domain.BranchInfos, args syncBranchStepsArgs) steps.List {
	list := steps.List{}
	for _, branch := range branchesToSync {
		syncBranchSteps(&list, branch, syncBranchStepsArgs{
			branches:           args.branches,
			branchTypes:        args.branches.Types,
			remotes:            args.remotes,
			isOffline:          args.isOffline,
			lineage:            args.lineage,
			mainBranch:         args.mainBranch,
			pullBranchStrategy: args.pullBranchStrategy,
			pushBranch:         args.pushBranch,
			pushHook:           args.pushHook,
			shouldSyncUpstream: args.shouldSyncUpstream,
			syncStrategy:       args.syncStrategy,
		})
	}
	list.Add(&step.CheckoutIfExists{Branch: args.branches.Initial})
	if args.remotes.HasOrigin() && args.shouldPushTags && !args.isOffline {
		list.Add(&step.PushTags{})
	}
	list.Wrap(steps.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: args.hasOpenChanges,
		MainBranch:       args.mainBranch,
		InitialBranch:    args.branches.Initial,
		PreviousBranch:   args.previousBranch,
	})
	return list
}

func syncBranchSteps(list *steps.List, branch domain.BranchInfo, args syncBranchStepsArgs) {
	if branch.SyncStatus == domain.SyncStatusDeletedAtRemote {
		syncDeletedBranchSteps(syncDeletedBranchArgs{
			branch:       branch,
			branchTypes:  args.branches.Types,
			lineage:      args.lineage,
			list:         list,
			mainBranch:   args.mainBranch,
			parent:       args.lineage.Parent(branch.LocalName),
			syncStrategy: args.syncStrategy,
		})
	} else {
		syncNonDeletedBranchSteps(list, branch, syncBranchStepsArgs{
			branchTypes:        args.branches.Types,
			remotes:            args.remotes,
			isOffline:          args.isOffline,
			lineage:            args.lineage,
			mainBranch:         args.mainBranch,
			pullBranchStrategy: args.pullBranchStrategy,
			pushBranch:         true,
			pushHook:           args.pushHook,
			shouldSyncUpstream: args.shouldSyncUpstream,
			syncStrategy:       args.syncStrategy,
		})
	}
}

type syncBranchStepsArgs struct {
	branches           domain.Branches
	branchTypes        domain.BranchTypes // TODO: branches already contains this
	remotes            domain.Remotes
	hasOpenChanges     bool
	hasUpstream        bool
	isOffline          bool
	lineage            config.Lineage
	mainBranch         domain.LocalBranchName
	previousBranch     domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	pushBranch         bool
	pushHook           bool
	shouldPushTags     bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

// syncDeletedBranchSteps provides a program that syncs a branch that was deleted at origin.
func syncDeletedBranchSteps(args syncDeletedBranchArgs) {
	if args.branchTypes.IsFeatureBranch(args.branch.LocalName) {
		syncDeleteFeatureBranchSteps(syncDeletedFeatureBranchArgs{
			branch:       args.branch,
			lineage:      args.lineage,
			list:         args.list,
			parent:       args.parent,
			syncStrategy: args.syncStrategy,
		})
	} else {
		syncDeletedPerennialBranchSteps(deletePerennialBranchStepsArgs{
			branch:     args.branch,
			lineage:    args.lineage,
			list:       args.list,
			mainBranch: args.mainBranch,
		})
	}
}

type syncDeletedBranchArgs struct {
	branch       domain.BranchInfo
	branchTypes  domain.BranchTypes
	lineage      config.Lineage
	list         *steps.List
	mainBranch   domain.LocalBranchName
	parent       domain.LocalBranchName
	syncStrategy config.SyncStrategy
}

// syncDeleteFeatureBranchSteps syncs a feare branch whose remote has been deleted.
// The parent branch must have been fully synced before calling this function.
func syncDeleteFeatureBranchSteps(args syncDeletedFeatureBranchArgs) {
	args.list.Add(&step.Checkout{Branch: args.branch.LocalName})
	pullParentBranchOfCurrentFeatureBranchStep(args.list, args.branch.LocalName, args.syncStrategy)
	// determine whether the now synced local branch still contains unshipped changes
	args.list.Add(&step.IfElse{
		Condition: func() (bool, error) {
			return args.backend.BranchHasUnmergedChanges(args.branch.LocalName, args.parent.Location())
		},
		TrueSteps: []step.Step{
			&step.QueueMessage{
				Message: fmt.Sprintf(messages.BranchDeletedHasUnmergedChanges, args.branch.LocalName),
			},
		},
		FalseSteps: []step.Step{
			&step.Checkout{Branch: args.parent},
			&step.DeleteLocalBranch{
				Branch: args.branch.LocalName,
				Force:  false,
				Parent: args.parent.Location(),
			},
			&step.RemoveBranchFromLineage{
				Branch: args.branch.LocalName,
			},
			&step.QueueMessage{
				Message: fmt.Sprintf(messages.BranchDeleted, args.branch.LocalName),
			},
		},
	})
}

type syncDeletedFeatureBranchArgs struct {
	backend      *git.BackendCommands
	branch       domain.BranchInfo
	lineage      config.Lineage
	list         *steps.List
	parent       domain.LocalBranchName
	syncStrategy config.SyncStrategy
}

func syncDeletedPerennialBranchSteps(args deletePerennialBranchStepsArgs) config.Lineage {
	result := removeBranchFromLineage(removeBranchFromLineageArgs{
		list:    args.list,
		branch:  args.branch.LocalName,
		parent:  args.mainBranch,
		lineage: args.lineage,
	})
	args.list.Add(&step.RemoveFromPerennialBranches{Branch: args.branch.LocalName})
	args.list.Add(&step.Checkout{Branch: args.mainBranch})
	args.list.Add(&step.DeleteLocalBranch{
		Branch: args.branch.LocalName,
		Force:  false,
		Parent: domain.Location(args.mainBranch),
	})
	args.list.Add(&step.QueueMessage{Message: fmt.Sprintf(messages.BranchDeleted, args.branch.LocalName)})
	return result
}

type deletePerennialBranchStepsArgs struct {
	list       *steps.List
	branch     domain.BranchInfo
	mainBranch domain.LocalBranchName
	lineage    config.Lineage
}

// syncBranchSteps provides the steps to sync a particular branch.
func syncNonDeletedBranchSteps(list *steps.List, branch domain.BranchInfo, args syncBranchStepsArgs) {
	isFeatureBranch := args.branchTypes.IsFeatureBranch(branch.LocalName)
	if !isFeatureBranch && !args.remotes.HasOrigin() {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	list.Add(&step.Checkout{Branch: branch.LocalName})
	if isFeatureBranch {
		syncFeatureBranchSteps(list, branch, args.syncStrategy)
	} else {
		syncPerennialBranchSteps(list, branch, args)
	}
	if args.pushBranch && args.remotes.HasOrigin() && !args.isOffline {
		switch {
		case !branch.HasTrackingBranch():
			list.Add(&step.CreateTrackingBranch{Branch: branch.LocalName, NoPushHook: !args.pushHook})
		case !isFeatureBranch:
			list.Add(&step.PushCurrentBranch{CurrentBranch: branch.LocalName, NoPushHook: !args.pushHook})
		default:
			pushFeatureBranchSteps(list, branch.LocalName, args.syncStrategy, args.pushHook)
		}
	}
}

// syncFeatureBranchSteps adds all the steps to sync the feature branch with the given name.
func syncFeatureBranchSteps(list *steps.List, branch domain.BranchInfo, syncStrategy config.SyncStrategy) {
	if branch.HasTrackingBranch() {
		pullTrackingBranchOfCurrentFeatureBranchStep(list, branch.RemoteName, syncStrategy)
	}
	pullParentBranchOfCurrentFeatureBranchStep(list, branch.LocalName, syncStrategy)
}

// syncPerennialBranchSteps adds all the steps to sync the perennial branch with the given name.
func syncPerennialBranchSteps(list *steps.List, branch domain.BranchInfo, args syncBranchStepsArgs) {
	if branch.HasTrackingBranch() {
		updateCurrentPerennialBranchStep(list, branch.RemoteName, args.pullBranchStrategy)
	}
	if branch.LocalName == args.mainBranch && args.hasUpstream && args.shouldSyncUpstream {
		list.Add(&step.FetchUpstream{Branch: args.mainBranch})
		list.Add(&step.RebaseBranch{Branch: domain.NewBranchName("upstream/" + args.mainBranch.String())})
	}
}

// pullTrackingBranchOfCurrentFeatureBranchStep adds the step to pull updates from the remote branch of the current feature branch into the current feature branch.
func pullTrackingBranchOfCurrentFeatureBranchStep(list *steps.List, trackingBranch domain.RemoteBranchName, strategy config.SyncStrategy) {
	switch strategy {
	case config.SyncStrategyMerge:
		list.Add(&step.Merge{Branch: trackingBranch.BranchName()})
	case config.SyncStrategyRebase:
		list.Add(&step.RebaseBranch{Branch: trackingBranch.BranchName()})
	}
}

// pullParentBranchOfCurrentFeatureBranchStep adds the step to pull updates from the parent branch of the current feature branch into the current feature branch.
func pullParentBranchOfCurrentFeatureBranchStep(list *steps.List, currentBranch domain.LocalBranchName, strategy config.SyncStrategy) {
	switch strategy {
	case config.SyncStrategyMerge:
		list.Add(&step.MergeParent{Branch: currentBranch})
	case config.SyncStrategyRebase:
		list.Add(&step.RebaseParent{Branch: currentBranch})
	}
}

// updateCurrentPerennialBranchStep provides the steps to update the current perennial branch with changes from the given other branch.
func updateCurrentPerennialBranchStep(list *steps.List, otherBranch domain.RemoteBranchName, strategy config.PullBranchStrategy) {
	switch strategy {
	case config.PullBranchStrategyMerge:
		list.Add(&step.Merge{Branch: otherBranch.BranchName()})
	case config.PullBranchStrategyRebase:
		list.Add(&step.RebaseBranch{Branch: otherBranch.BranchName()})
	}
}

func pushFeatureBranchSteps(list *steps.List, branch domain.LocalBranchName, syncStrategy config.SyncStrategy, pushHook bool) {
	switch syncStrategy {
	case config.SyncStrategyMerge:
		list.Add(&step.PushCurrentBranch{CurrentBranch: branch, NoPushHook: !pushHook})
	case config.SyncStrategyRebase:
		list.Add(&step.ForcePushCurrentBranch{NoPushHook: !pushHook})
	}
}
