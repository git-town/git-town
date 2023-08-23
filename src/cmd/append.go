package cmd

import (
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
	repo, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, exit, err := determineAppendConfig(domain.NewLocalBranchName(arg), &repo)
	if err != nil || exit {
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
	branches            domain.Branches
	branchesToSync      domain.BranchInfos
	hasOpenChanges      bool
	remotes             domain.Remotes
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

func determineAppendConfig(targetBranch domain.LocalBranchName, repo *execute.RepoData) (*appendConfig, bool, error) {
	branches, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	fc := failure.Collector{}
	remotes := fc.Remotes(repo.Runner.Backend.Remotes())
	mainBranch := repo.Runner.Config.MainBranch()
	pushHook := fc.Bool(repo.Runner.Config.PushHook())
	pullBranchStrategy := fc.PullBranchStrategy(repo.Runner.Config.PullBranchStrategy())
	hasOpenChanges := fc.Bool(repo.Runner.Backend.HasOpenChanges())
	shouldNewBranchPush := fc.Bool(repo.Runner.Config.ShouldNewBranchPush())
	if fc.Err != nil {
		return nil, false, fc.Err
	}
	if branches.All.HasLocalBranch(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branches.All.HasMatchingRemoteBranchFor(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	lineage := repo.Runner.Config.Lineage()
	updated, err := validate.KnowsBranchAncestors(branches.Initial, validate.KnowsBranchAncestorsArgs{
		DefaultBranch: mainBranch,
		Backend:       &repo.Runner.Backend,
		AllBranches:   branches.All,
		Lineage:       lineage,
		BranchTypes:   branches.Types,
		MainBranch:    mainBranch,
	})
	if err != nil {
		return nil, false, err
	}
	if updated {
		lineage = repo.Runner.Config.Lineage() // refresh lineage after ancestry changes
	}
	branchNamesToSync := lineage.BranchAndAncestors(branches.Initial)
	branchesToSync := fc.BranchesSyncStatus(branches.All.Select(branchNamesToSync))
	syncStrategy := fc.SyncStrategy(repo.Runner.Config.SyncStrategy())
	shouldSyncUpstream := fc.Bool(repo.Runner.Config.ShouldSyncUpstream())
	return &appendConfig{
		branches:            branches,
		branchesToSync:      branchesToSync,
		hasOpenChanges:      hasOpenChanges,
		remotes:             remotes,
		isOffline:           repo.IsOffline,
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
	}, false, fc.Err
}

func appendStepList(config *appendConfig) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.branchesToSync {
		syncBranchSteps(&list, syncBranchStepsArgs{
			branch:             branch,
			branchTypes:        config.branches.Types,
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
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch.Location()})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
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
