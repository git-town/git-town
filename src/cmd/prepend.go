package cmd

import (
	"fmt"

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

const prependDesc = "Creates a new feature branch as the parent of the current branch"

const prependHelp = `
Syncs the parent branch,
cuts a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
pushes the new feature branch to the origin repository
(if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for upstream remote options.
`

func prependCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "prepend <branch>",
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		Short:   prependDesc,
		Long:    long(prependDesc, prependHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return prepend(args, readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func prepend(args []string, debug bool) error {
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
	config, err := determinePrependConfig(args, &repo.Runner, repo.IsOffline)
	if err != nil {
		return err
	}
	stepList, err := prependStepList(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "prepend",
		RunStepList: stepList,
	}
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &runState,
		Run:       &repo.Runner,
		Connector: nil,
		RootDir:   repo.RootDir,
	})
}

type prependConfig struct {
	branchDurations     config.BranchDurations
	branchesToSync      git.BranchesSyncStatus
	hasOpenChanges      bool
	remotes             config.Remotes
	initialBranch       string
	isOffline           bool
	lineage             config.Lineage
	mainBranch          string
	previousBranch      string
	pullBranchStrategy  config.PullBranchStrategy
	pushHook            bool
	parentBranch        string
	shouldSyncUpstream  bool
	shouldNewBranchPush bool
	syncStrategy        config.SyncStrategy
	targetBranch        string
}

func determinePrependConfig(args []string, run *git.ProdRunner, isOffline bool) (*prependConfig, error) {
	fc := failure.Collector{}
	branches := fc.Branches(execute.LoadBranches(run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	}))
	previousBranch := run.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges := fc.Bool(run.Backend.HasOpenChanges())
	remotes := fc.Strings(run.Backend.Remotes())
	shouldNewBranchPush := fc.Bool(run.Config.ShouldNewBranchPush())
	pushHook := fc.Bool(run.Config.PushHook())
	mainBranch := run.Config.MainBranch()
	syncStrategy := fc.SyncStrategy(run.Config.SyncStrategy())
	pullBranchStrategy := fc.PullBranchStrategy(run.Config.PullBranchStrategy())
	shouldSyncUpstream := fc.Bool(run.Config.ShouldSyncUpstream())
	// TODO: use fc all the way to the end
	if fc.Err != nil {
		return nil, fc.Err
	}
	targetBranch := args[0]
	if branches.All.Contains(targetBranch) {
		return nil, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if !branches.Durations.IsFeatureBranch(branches.Initial) {
		return nil, fmt.Errorf(messages.SetParentNoFeatureBranch, branches.Initial)
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
		lineage = run.Config.Lineage()
	}
	branchNamesToSync := lineage.BranchAndAncestors(branches.Initial)
	branchesToSync, err := branches.All.Select(branchNamesToSync)
	return &prependConfig{
		branchDurations:     branches.Durations,
		branchesToSync:      branchesToSync,
		hasOpenChanges:      hasOpenChanges,
		remotes:             remotes,
		initialBranch:       branches.Initial,
		isOffline:           isOffline,
		lineage:             lineage,
		mainBranch:          mainBranch,
		previousBranch:      previousBranch,
		pullBranchStrategy:  pullBranchStrategy,
		pushHook:            pushHook,
		parentBranch:        lineage.Parent(branches.Initial),
		shouldNewBranchPush: shouldNewBranchPush,
		shouldSyncUpstream:  shouldSyncUpstream,
		syncStrategy:        syncStrategy,
		targetBranch:        targetBranch,
	}, err
}

func prependStepList(config *prependConfig) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branchToSync := range config.branchesToSync {
		syncBranchSteps(&list, syncBranchStepsArgs{
			branch:             branchToSync,
			branchDurations:    config.branchDurations,
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
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.initialBranch, ParentBranch: config.targetBranch})
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
