package cmd

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
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
	config, err := determineAppendConfig(domain.NewLocalBranchName(arg), &repo.Runner, repo.IsOffline)
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
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &runState,
		Run:       &repo.Runner,
		Connector: nil,
		RootDir:   repo.RootDir,
	})
}

type appendConfig struct {
	durations           config.BranchDurations
	branchesToSync      git.BranchesSyncStatus
	hasOpenChanges      bool
	remotes             config.Remotes
	initialBranch       domain.LocalBranchName
	isOffline           bool
	lineage             config.Lineage
	mainBranch          domain.LocalBranchName
	pushHook            bool
	parentBranch        domain.LocalBranchName
	previousBranch      domain.LocalBranchName
	pullBranchStrategy  config.PullBranchStrategy
	shouldNewBranchPush bool
	shouldSyncUpstream  bool
	syncStrategy        config.SyncStrategy
	targetBranch        domain.LocalBranchName
}

func determineAppendConfig(targetBranch domain.LocalBranchName, run *git.ProdRunner, isOffline bool) (*appendConfig, error) {
	branches, err := execute.LoadBranches(run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return nil, err
	}
	previousBranch := run.Backend.PreviouslyCheckedOutBranch()
	fc := failure.Collector{}
	remotes := fc.Strings(run.Backend.Remotes())
	mainBranch := run.Config.MainBranch()
	pushHook := fc.Bool(run.Config.PushHook())
	pullBranchStrategy := fc.PullBranchStrategy(run.Config.PullBranchStrategy())
	hasOpenChanges := fc.Bool(run.Backend.HasOpenChanges())
	shouldNewBranchPush := fc.Bool(run.Config.ShouldNewBranchPush())
	if fc.Err != nil {
		return nil, fc.Err
	}
	if branches.All.ContainsLocalBranch(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branches.All.KnowsRemoteBranch(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	lineage := run.Config.Lineage()
	updated, err := validate.KnowsBranchAncestors(branches.Initial, validate.KnowsBranchAncestorsArgs{
		DefaultBranch:   mainBranch,
		Backend:         &run.Backend,
		AllBranches:     branches.All,
		Lineage:         lineage,
		BranchDurations: branches.Durations,
		MainBranch:      mainBranch,
	})
	if err != nil {
		return nil, err
	}
	if updated {
		lineage = run.Config.Lineage() // refresh lineage after ancestry changes
	}
	branchNamesToSync := lineage.BranchAndAncestors(branches.Initial)
	branchesToSync := fc.BranchesSyncStatus(branches.All.Select(branchNamesToSync))
	syncStrategy := fc.SyncStrategy(run.Config.SyncStrategy())
	shouldSyncUpstream := fc.Bool(run.Config.ShouldSyncUpstream())
	return &appendConfig{
		durations:           branches.Durations,
		branchesToSync:      branchesToSync,
		hasOpenChanges:      hasOpenChanges,
		remotes:             remotes,
		initialBranch:       branches.Initial,
		isOffline:           isOffline,
		lineage:             lineage,
		mainBranch:          mainBranch,
		pushHook:            pushHook,
		parentBranch:        branches.Initial,
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
		syncBranchSteps(&list, syncBranchStepsArgs{
			branch:             branch,
			branchDurations:    config.durations,
			isOffline:          config.isOffline,
			lineage:            config.lineage,
			remotes:            config.remotes,
			mainBranch:         config.mainBranch,
			pullBranchStrategy: config.pullBranchStrategy,
			pushBranch:         true,
			pushHook:           config.pushHook,
			shouldSyncUpstream: config.shouldSyncUpstream,
			syncStrategy:       config.syncStrategy,
		})
	}
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch.Location})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	list.Add(&steps.CheckoutStep{Branch: config.targetBranch})
	if config.remotes.HasOrigin() && config.shouldNewBranchPush && !config.isOffline {
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
