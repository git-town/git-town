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
	stepList, err := appendStepList(config, &repo.Runner)
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
	branchesToSync      git.BranchesSyncStatus
	hasOrigin           bool
	isOffline           bool
	mainBranch          string
	noPushHook          bool
	parentBranch        string
	shouldNewBranchPush bool
	syncStrategy        config.SyncStrategy
	targetBranch        string
}

func determineAppendConfig(targetBranch string, run *git.ProdRunner, allBranches git.BranchesSyncStatus, currentBranchName string, isOffline bool) (*appendConfig, error) {
	fc := failure.Collector{}
	hasOrigin := fc.Bool(run.Backend.HasOrigin())
	mainBranch := run.Config.MainBranch()
	pushHook := fc.Bool(run.Config.PushHook())
	shouldNewBranchPush := fc.Bool(run.Config.ShouldNewBranchPush())
	if fc.Err != nil {
		return nil, fc.Err
	}
	if allBranches.Contains(targetBranch) {
		fc.Fail(messages.BranchAlreadyExists, targetBranch)
	}
	fc.Check(validate.KnowsBranchAncestors(currentBranchName, mainBranch, &run.Backend))
	lineage := run.Config.Lineage()
	branchNamesToSync := lineage.BranchAndAncestors(currentBranchName)
	branchesToSync := fc.BranchesSyncStatus(allBranches.Select(branchNamesToSync))
	syncStrategy := fc.SyncStrategy(run.Config.SyncStrategy())
	return &appendConfig{
		branchesToSync:      branchesToSync,
		isOffline:           isOffline,
		hasOrigin:           hasOrigin,
		mainBranch:          mainBranch,
		noPushHook:          !pushHook,
		parentBranch:        currentBranchName,
		shouldNewBranchPush: shouldNewBranchPush,
		syncStrategy:        syncStrategy,
		targetBranch:        targetBranch,
	}, fc.Err
}

func appendStepList(config *appendConfig, run *git.ProdRunner) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.branchesToSync {
		updateBranchSteps(&list, updateBranchStepsArgs{
			branch:       branch,
			isOffline:    config.isOffline,
			hasOrigin:    config.hasOrigin,
			mainBranch:   config.mainBranch,
			pushBranch:   true,
			run:          run,
			syncStrategy: config.syncStrategy,
		})
	}
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	list.Add(&steps.CheckoutStep{Branch: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		list.Add(&steps.CreateTrackingBranchStep{Branch: config.targetBranch, NoPushHook: config.noPushHook})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, &run.Backend, config.mainBranch)
	return list.Result()
}
