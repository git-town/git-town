package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/flags"
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
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, exit, err := determinePrependConfig(args, &repo)
	if err != nil || exit {
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
	branches            domain.Branches
	branchesToSync      domain.BranchInfos
	hasOpenChanges      bool
	remotes             domain.Remotes
	isOffline           bool
	lineage             config.Lineage
	mainBranch          domain.LocalBranchName
	previousBranch      domain.LocalBranchName
	pullBranchStrategy  config.PullBranchStrategy
	pushHook            bool
	parentBranch        domain.LocalBranchName
	shouldSyncUpstream  bool
	shouldNewBranchPush bool
	syncStrategy        config.SyncStrategy
	targetBranch        domain.LocalBranchName
}

func determinePrependConfig(args []string, repo *execute.OpenRepoResult) (*prependConfig, bool, error) {
	branches, _, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	fc := failure.Collector{}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges := fc.Bool(repo.Runner.Backend.HasOpenChanges())
	remotes := fc.Remotes(repo.Runner.Backend.Remotes())
	shouldNewBranchPush := fc.Bool(repo.Runner.Config.ShouldNewBranchPush())
	pushHook := fc.Bool(repo.Runner.Config.PushHook())
	mainBranch := repo.Runner.Config.MainBranch()
	syncStrategy := fc.SyncStrategy(repo.Runner.Config.SyncStrategy())
	pullBranchStrategy := fc.PullBranchStrategy(repo.Runner.Config.PullBranchStrategy())
	shouldSyncUpstream := fc.Bool(repo.Runner.Config.ShouldSyncUpstream())
	targetBranch := domain.NewLocalBranchName(args[0])
	if branches.All.HasLocalBranch(targetBranch) {
		return nil, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branches.All.HasMatchingRemoteBranchFor(targetBranch) {
		return nil, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	if !branches.Types.IsFeatureBranch(branches.Initial) {
		return nil, false, fmt.Errorf(messages.SetParentNoFeatureBranch, branches.Initial)
	}
	lineage := repo.Runner.Config.Lineage()
	updated := fc.Bool(validate.KnowsBranchAncestors(branches.Initial, validate.KnowsBranchAncestorsArgs{
		DefaultBranch: mainBranch,
		Backend:       &repo.Runner.Backend,
		AllBranches:   branches.All,
		Lineage:       lineage,
		BranchTypes:   branches.Types,
		MainBranch:    mainBranch,
	}))
	if updated {
		lineage = repo.Runner.Config.Lineage()
	}
	branchNamesToSync := lineage.BranchAndAncestors(branches.Initial)
	branchesToSync := fc.BranchesSyncStatus(branches.All.Select(branchNamesToSync))
	return &prependConfig{
		branches:            branches,
		branchesToSync:      branchesToSync,
		hasOpenChanges:      hasOpenChanges,
		remotes:             remotes,
		isOffline:           repo.IsOffline,
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
	}, false, fc.Err
}

func prependStepList(config *prependConfig) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branchToSync := range config.branchesToSync {
		syncBranchSteps(&list, syncBranchStepsArgs{
			branch:             branchToSync,
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
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch.Location()})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.branches.Initial, ParentBranch: config.targetBranch})
	list.Add(&steps.CheckoutStep{Branch: config.targetBranch})
	if config.remotes.HasOrigin() && config.shouldNewBranchPush && !config.isOffline {
		list.Add(&steps.CreateTrackingBranchStep{Branch: config.targetBranch, NoPushHook: !config.pushHook})
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
