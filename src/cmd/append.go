package cmd

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/validate"
	"github.com/spf13/cobra"
)

const appendDesc = "Creates a new feature branch as a child of the current branch"

const appendHelp = `
Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the origin repository
(if and only if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`

func appendCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "append <branch>",
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		Short:   appendDesc,
		Long:    long(appendDesc, appendHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAppend(args[0], readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func runAppend(arg string, debug bool) error {
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
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
	allBranches, currentBranch, err := execute.LoadBranches(&repo.Runner, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	config, err := determineAppendConfig(arg, &repo.Runner, allBranches, currentBranch, repo.IsOffline)
	if err != nil {
		return err
	}
	stepList, err := appendStepList(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "append",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &repo.Runner, nil, repo.RootDir)
}

type appendConfig struct {
	branchDurations     config.BranchDurations
	branchesToSync      git.BranchesSyncStatus
	hasOpenChanges      bool
	hasOrigin           bool
	hasUpstream         bool
	initialBranch       string
	isOffline           bool
	lineage             config.Lineage
	mainBranch          string
	pushHook            bool
	parentBranch        string
	previousBranch      string
	pullBranchStrategy  config.PullBranchStrategy
	shouldNewBranchPush bool
	shouldSyncUpstream  bool
	syncStrategy        config.SyncStrategy
	targetBranch        string
}

func determineAppendConfig(targetBranch string, run *git.ProdRunner, allBranches git.BranchesSyncStatus, initialBranch string, isOffline bool) (*appendConfig, error) {
	previousBranch := run.Backend.PreviouslyCheckedOutBranch()
	fc := failure.Collector{}
	hasOrigin := fc.Bool(run.Backend.HasOrigin())
	mainBranch := run.Config.MainBranch()
	pushHook := fc.Bool(run.Config.PushHook())
	hasUpstream := fc.Bool(run.Backend.HasUpstream())
	pullBranchStrategy := fc.PullBranchStrategy(run.Config.PullBranchStrategy())
	hasOpenChanges := fc.Bool(run.Backend.HasOpenChanges())
	shouldNewBranchPush := fc.Bool(run.Config.ShouldNewBranchPush())
	if fc.Err != nil {
		return nil, fc.Err
	}
	if allBranches.Contains(targetBranch) {
		fc.Fail(messages.BranchAlreadyExists, targetBranch)
	}
	branchDurations := run.Config.BranchDurations()
	fc.Check(validate.KnowsBranchAncestors(initialBranch, validate.KnowsBranchAncestorsArgs{
		DefaultBranch:   mainBranch,
		Backend:         &run.Backend,
		AllBranches:     allBranches,
		Lineage:         run.Config.Lineage(),
		BranchDurations: branchDurations,
		MainBranch:      mainBranch,
	}))
	lineage := run.Config.Lineage() // refresh lineage after ancestry changes
	branchNamesToSync := lineage.BranchAndAncestors(initialBranch)
	branchesToSync := fc.BranchesSyncStatus(allBranches.Select(branchNamesToSync))
	syncStrategy := fc.SyncStrategy(run.Config.SyncStrategy())
	shouldSyncUpstream := fc.Bool(run.Config.ShouldSyncUpstream())
	return &appendConfig{
		branchDurations:     branchDurations,
		branchesToSync:      branchesToSync,
		hasOpenChanges:      hasOpenChanges,
		hasOrigin:           hasOrigin,
		hasUpstream:         hasUpstream,
		initialBranch:       initialBranch,
		isOffline:           isOffline,
		lineage:             lineage,
		mainBranch:          mainBranch,
		pushHook:            pushHook,
		parentBranch:        initialBranch,
		previousBranch:      previousBranch,
		pullBranchStrategy:  pullBranchStrategy,
		shouldNewBranchPush: shouldNewBranchPush,
		shouldSyncUpstream:  shouldSyncUpstream,
		syncStrategy:        syncStrategy,
		targetBranch:        targetBranch,
	}, fc.Err
}

func appendStepList(config *appendConfig) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.branchesToSync {
		updateBranchSteps(&list, updateBranchStepsArgs{
			branch:             branch,
			branchDurations:    config.branchDurations,
			isOffline:          config.isOffline,
			lineage:            config.lineage,
			hasOrigin:          config.hasOrigin,
			hasUpstream:        config.hasUpstream,
			mainBranch:         config.mainBranch,
			pullBranchStrategy: config.pullBranchStrategy,
			pushBranch:         true,
			pushHook:           config.pushHook,
			shouldSyncUpstream: config.shouldSyncUpstream,
			syncStrategy:       config.syncStrategy,
		})
	}
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	list.Add(&steps.CheckoutStep{Branch: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		list.Add(&steps.CreateTrackingBranchStep{Branch: config.targetBranch, NoPushHook: !config.pushHook})
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
