package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
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
	repo, commandsCounter, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           dryRun,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineSyncConfig(all, repo)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "sync",
		CommandsRun:         commandsCounter,
		InitialActiveBranch: initialBranchesSnapshot.Active,
		MessagesToUser:      runstate.NewMessages(),
		RunSteps:            syncBranchesSteps(config),
	}
	return runvm.Execute(runvm.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               nil,
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
	pushHook           bool
	shouldPushTags     bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

func determineSyncConfig(allFlag bool, repo *execute.OpenRepoResult) (*syncConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage()
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
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
func syncBranchesSteps(config *syncConfig) steps.List {
	list := steps.List{}
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
	list.Add(&step.Checkout{Branch: config.branches.Initial})
	if config.remotes.HasOrigin() && config.shouldPushTags && !config.isOffline {
		list.Add(&step.PushTags{})
	}
	list.Wrap(steps.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return list
}

// syncBranchSteps provides the steps to sync a particular branch.
func syncBranchSteps(list *steps.List, args syncBranchStepsArgs) {
	isFeatureBranch := args.branchTypes.IsFeatureBranch(args.branch.LocalName)
	if !isFeatureBranch && !args.remotes.HasOrigin() {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	list.Add(&step.Checkout{Branch: args.branch.LocalName})
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
			list.Add(&step.CreateTrackingBranch{Branch: args.branch.LocalName, NoPushHook: !args.pushHook})
		case !isFeatureBranch:
			list.Add(&step.PushCurrentBranch{CurrentBranch: args.branch.LocalName, NoPushHook: !args.pushHook})
		default:
			pushFeatureBranchSteps(list, args.branch.LocalName, args.syncStrategy, args.pushHook)
		}
	}
}

type syncBranchStepsArgs struct {
	branch             domain.BranchInfo
	branchTypes        domain.BranchTypes
	remotes            domain.Remotes
	isOffline          bool
	lineage            config.Lineage
	mainBranch         domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	pushBranch         bool
	pushHook           bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

// syncFeatureBranchSteps adds all the steps to sync the feature branch with the given name.
func syncFeatureBranchSteps(list *steps.List, branch domain.BranchInfo, lineage config.Lineage, syncStrategy config.SyncStrategy) {
	if branch.HasTrackingBranch() {
		pullTrackingBranchOfCurrentFeatureBranchStep(list, branch.RemoteName, syncStrategy)
	}
	pullParentBranchOfCurrentFeatureBranchStep(list, lineage.Parent(branch.LocalName), syncStrategy)
}

// syncPerennialBranchSteps adds all the steps to sync the perennial branch with the given name.
func syncPerennialBranchSteps(list *steps.List, args syncPerennialBranchStepsArgs) {
	if args.branch.HasTrackingBranch() {
		updateCurrentPerennialBranchStep(list, args.branch.RemoteName, args.pullBranchStrategy)
	}
	if args.branch.LocalName == args.mainBranch && args.hasUpstream && args.shouldSyncUpstream {
		list.Add(&step.FetchUpstream{Branch: args.mainBranch})
		list.Add(&step.RebaseBranch{Branch: domain.NewBranchName("upstream/" + args.mainBranch.String())})
	}
}

type syncPerennialBranchStepsArgs struct {
	branch             domain.BranchInfo
	mainBranch         domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	shouldSyncUpstream bool
	hasUpstream        bool
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
func pullParentBranchOfCurrentFeatureBranchStep(list *steps.List, parentBranch domain.LocalBranchName, strategy config.SyncStrategy) {
	switch strategy {
	case config.SyncStrategyMerge:
		list.Add(&step.Merge{Branch: parentBranch.BranchName()})
	case config.SyncStrategyRebase:
		list.Add(&step.RebaseBranch{Branch: parentBranch.BranchName()})
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
