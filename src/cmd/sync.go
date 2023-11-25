package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/vm/interpreter"
	"github.com/git-town/git-town/v10/src/vm/opcode"
	"github.com/git-town/git-town/v10/src/vm/program"
	"github.com/git-town/git-town/v10/src/vm/runstate"
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
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addAllFlag, readAllFlag := flags.Bool("all", "a", "Sync all local branches", flags.FlagTypeNonPersistent)
	cmd := cobra.Command{
		Use:     "sync",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   syncDesc,
		Long:    long(syncDesc, fmt.Sprintf(syncHelp, config.KeySyncUpstream)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSync(readAllFlag(cmd), readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addAllFlag(&cmd)
	addVerboseFlag(&cmd)
	addDryRunFlag(&cmd)
	return &cmd
}

func executeSync(all, dryRun, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineSyncConfig(all, repo, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := program.Program{}
	syncBranchesProgram(syncBranchesProgramArgs{
		syncBranchProgramArgs: syncBranchProgramArgs{
			branchTypes:        config.branches.Types,
			remotes:            config.remotes,
			isOffline:          config.isOffline,
			lineage:            config.lineage,
			program:            &runProgram,
			mainBranch:         config.mainBranch,
			pullBranchStrategy: config.pullBranchStrategy,
			pushBranch:         true,
			pushHook:           config.pushHook,
			shouldSyncUpstream: config.shouldSyncUpstream,
			syncStrategy:       config.syncStrategy,
		},
		branchesToSync: config.branchesToSync,
		hasOpenChanges: config.hasOpenChanges,
		initialBranch:  config.branches.Initial,
		previousBranch: config.previousBranch,
		shouldPushTags: config.shouldPushTags,
	})
	runState := runstate.RunState{
		Command:             "sync",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          runProgram,
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               nil,
		Verbose:                 verbose,
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
	isOffline          bool
	lineage            config.Lineage
	mainBranch         domain.LocalBranchName
	previousBranch     domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	pushHook           bool
	remotes            domain.Remotes
	shouldPushTags     bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

func determineSyncConfig(allFlag bool, repo *execute.OpenRepoResult, verbose bool) (*syncConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
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
	if allFlag {
		localBranches := branches.All.LocalBranches()
		branches.Types, lineage, err = execute.EnsureKnownBranchesAncestry(execute.EnsureKnownBranchesAncestryArgs{
			AllBranches: localBranches,
			BranchTypes: branches.Types,
			Lineage:     lineage,
			MainBranch:  mainBranch,
			Runner:      &repo.Runner,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSnapshot, false, err
		}
		branchNamesToSync = localBranches.Names()
		shouldPushTags = true
	} else {
		branches.Types, lineage, err = execute.EnsureKnownBranchAncestry(branches.Initial, execute.EnsureKnownBranchAncestryArgs{
			AllBranches:   branches.All,
			BranchTypes:   branches.Types,
			DefaultBranch: mainBranch,
			Lineage:       lineage,
			MainBranch:    mainBranch,
			Runner:        &repo.Runner,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSnapshot, false, err
		}
	}
	if !allFlag {
		branchNamesToSync = domain.LocalBranchNames{branches.Initial}
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

// syncBranchesProgram provides the program for the "git sync" command.
func syncBranchesProgram(args syncBranchesProgramArgs) {
	for _, branch := range args.branchesToSync {
		syncBranchProgram(branch, args.syncBranchProgramArgs)
	}
	args.program.Add(&opcode.CheckoutIfExists{Branch: args.initialBranch})
	if args.remotes.HasOrigin() && args.shouldPushTags && !args.isOffline {
		args.program.Add(&opcode.PushTags{})
	}
	wrap(args.program, wrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: args.hasOpenChanges,
		MainBranch:       args.mainBranch,
		InitialBranch:    args.initialBranch,
		PreviousBranch:   args.previousBranch,
	})
}

type syncBranchesProgramArgs struct {
	syncBranchProgramArgs
	branchesToSync domain.BranchInfos
	hasOpenChanges bool
	initialBranch  domain.LocalBranchName
	previousBranch domain.LocalBranchName
	shouldPushTags bool
}

func syncBranchProgram(branch domain.BranchInfo, args syncBranchProgramArgs) {
	if branch.SyncStatus == domain.SyncStatusDeletedAtRemote {
		syncDeletedBranchProgram(args.program, branch, args)
	} else {
		syncNonDeletedBranchProgram(args.program, branch, args)
	}
}

type syncBranchProgramArgs struct {
	branchTypes        domain.BranchTypes
	isOffline          bool
	lineage            config.Lineage
	program            *program.Program
	mainBranch         domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	pushBranch         bool
	pushHook           bool
	remotes            domain.Remotes
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

// syncDeletedBranchProgram provides a program that syncs a branch that was deleted at origin.
func syncDeletedBranchProgram(list *program.Program, branch domain.BranchInfo, args syncBranchProgramArgs) {
	if args.branchTypes.IsFeatureBranch(branch.LocalName) {
		syncDeletedFeatureBranchProgram(list, branch, args)
	} else {
		syncDeletedPerennialBranchProgram(list, branch, args)
	}
}

// syncDeletedFeatureBranchProgram syncs a feare branch whose remote has been deleted.
// The parent branch must have been fully synced before calling this function.
func syncDeletedFeatureBranchProgram(list *program.Program, branch domain.BranchInfo, args syncBranchProgramArgs) {
	list.Add(&opcode.Checkout{Branch: branch.LocalName})
	pullParentBranchOfCurrentFeatureBranchOpcode(list, branch.LocalName, args.syncStrategy)
	list.Add(&opcode.DeleteBranchIfEmptyAtRuntime{Branch: branch.LocalName})
}

func syncDeletedPerennialBranchProgram(list *program.Program, branch domain.BranchInfo, args syncBranchProgramArgs) {
	removeBranchFromLineage(removeBranchFromLineageArgs{
		program: list,
		branch:  branch.LocalName,
		parent:  args.mainBranch,
		lineage: args.lineage,
	})
	list.Add(&opcode.RemoveFromPerennialBranches{Branch: branch.LocalName})
	list.Add(&opcode.Checkout{Branch: args.mainBranch})
	list.Add(&opcode.DeleteLocalBranch{
		Branch: branch.LocalName,
		Force:  false,
	})
	list.Add(&opcode.QueueMessage{Message: fmt.Sprintf(messages.BranchDeleted, branch.LocalName)})
}

// syncNonDeletedBranchProgram provides the opcode to sync a particular branch.
func syncNonDeletedBranchProgram(list *program.Program, branch domain.BranchInfo, args syncBranchProgramArgs) {
	isFeatureBranch := args.branchTypes.IsFeatureBranch(branch.LocalName)
	if !isFeatureBranch && !args.remotes.HasOrigin() {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	list.Add(&opcode.Checkout{Branch: branch.LocalName})
	if isFeatureBranch {
		syncFeatureBranchProgram(list, branch, args.syncStrategy)
	} else {
		syncPerennialBranchProgram(branch, args)
	}
	if args.pushBranch && args.remotes.HasOrigin() && !args.isOffline {
		switch {
		case !branch.HasTrackingBranch():
			list.Add(&opcode.CreateTrackingBranch{Branch: branch.LocalName, NoPushHook: !args.pushHook})
		case !isFeatureBranch:
			list.Add(&opcode.PushCurrentBranch{CurrentBranch: branch.LocalName, NoPushHook: !args.pushHook})
		default:
			pushFeatureBranchProgram(list, branch.LocalName, args.syncStrategy, args.pushHook)
		}
	}
}

// syncFeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func syncFeatureBranchProgram(list *program.Program, branch domain.BranchInfo, syncStrategy config.SyncStrategy) {
	if branch.HasTrackingBranch() {
		pullTrackingBranchOfCurrentFeatureBranchOpcode(list, branch.RemoteName, syncStrategy)
	}
	pullParentBranchOfCurrentFeatureBranchOpcode(list, branch.LocalName, syncStrategy)
}

// syncPerennialBranchProgram adds the opcodes to sync the perennial branch with the given name.
func syncPerennialBranchProgram(branch domain.BranchInfo, args syncBranchProgramArgs) {
	if branch.HasTrackingBranch() {
		updateCurrentPerennialBranchOpcode(args.program, branch.RemoteName, args.pullBranchStrategy)
	}
	if branch.LocalName == args.mainBranch && args.remotes.HasUpstream() && args.shouldSyncUpstream {
		args.program.Add(&opcode.FetchUpstream{Branch: args.mainBranch})
		args.program.Add(&opcode.RebaseBranch{Branch: domain.NewBranchName("upstream/" + args.mainBranch.String())})
	}
}

// pullTrackingBranchOfCurrentFeatureBranchOpcode adds the opcode to pull updates from the remote branch of the current feature branch into the current feature branch.
func pullTrackingBranchOfCurrentFeatureBranchOpcode(list *program.Program, trackingBranch domain.RemoteBranchName, strategy config.SyncStrategy) {
	switch strategy {
	case config.SyncStrategyMerge:
		list.Add(&opcode.Merge{Branch: trackingBranch.BranchName()})
	case config.SyncStrategyRebase:
		list.Add(&opcode.RebaseBranch{Branch: trackingBranch.BranchName()})
	}
}

// pullParentBranchOfCurrentFeatureBranchOpcode adds the opcode to pull updates from the parent branch of the current feature branch into the current feature branch.
func pullParentBranchOfCurrentFeatureBranchOpcode(list *program.Program, currentBranch domain.LocalBranchName, strategy config.SyncStrategy) {
	switch strategy {
	case config.SyncStrategyMerge:
		list.Add(&opcode.MergeParent{CurrentBranch: currentBranch})
	case config.SyncStrategyRebase:
		list.Add(&opcode.RebaseParent{CurrentBranch: currentBranch})
	}
}

// updateCurrentPerennialBranchOpcode provides the opcode to update the current perennial branch with changes from the given other branch.
func updateCurrentPerennialBranchOpcode(list *program.Program, otherBranch domain.RemoteBranchName, strategy config.PullBranchStrategy) {
	switch strategy {
	case config.PullBranchStrategyMerge:
		list.Add(&opcode.Merge{Branch: otherBranch.BranchName()})
	case config.PullBranchStrategyRebase:
		list.Add(&opcode.RebaseBranch{Branch: otherBranch.BranchName()})
	}
}

func pushFeatureBranchProgram(list *program.Program, branch domain.LocalBranchName, syncStrategy config.SyncStrategy, pushHook bool) {
	switch syncStrategy {
	case config.SyncStrategyMerge:
		list.Add(&opcode.PushCurrentBranch{CurrentBranch: branch, NoPushHook: !pushHook})
	case config.SyncStrategyRebase:
		list.Add(&opcode.ForcePushCurrentBranch{NoPushHook: !pushHook})
	}
}
