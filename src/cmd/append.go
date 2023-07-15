package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
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
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: true,
		ValidateGitversion:    true,
		ValidateIsRepository:  true,
		ValidateIsConfigured:  true,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineAppendConfig(arg, &run)
	if err != nil {
		return err
	}
	stepList, err := appendStepList(config, &run)
	if err != nil {
		return err
	}
	runState := runstate.New("append", stepList)
	return runstate.Execute(runState, &run, nil)
}

type appendConfig struct {
	ancestorBranches    config.Lineage
	hasOrigin           bool
	isOffline           bool
	mainBranch          string
	noPushHook          bool
	parentBranch        config.BranchWithParent
	shouldNewBranchPush bool
	targetBranch        string
}

func determineAppendConfig(targetBranch string, run *git.ProdRunner) (*appendConfig, error) {
	fc := failure.Collector{}
	lineage := run.Config.Lineage()
	parentBranchName := fc.String(run.Backend.CurrentBranch())
	parentBranch := lineage.Lookup(parentBranchName)
	if parentBranch == nil {
		return nil, fmt.Errorf("cannot find current branch %q in lineage", parentBranchName)
	}
	hasOrigin := fc.Bool(run.Backend.HasOrigin())
	isOffline := fc.Bool(run.Config.IsOffline())
	mainBranch := run.Config.MainBranch()
	pushHook := fc.Bool(run.Config.PushHook())
	shouldNewBranchPush := fc.Bool(run.Config.ShouldNewBranchPush())
	if fc.Err != nil {
		return nil, fc.Err
	}
	if hasOrigin && !isOffline {
		fc.Check(run.Frontend.Fetch())
	}
	hasTargetBranch := fc.Bool(run.Backend.HasLocalOrOriginBranch(targetBranch, mainBranch))
	if hasTargetBranch {
		fc.Fail("a branch named %q already exists", targetBranch)
	}
	fc.Check(validate.KnowsBranchAncestors(parentBranchName, run.Config.MainBranch(), &run.Backend))
	ancestorBranches := lineage.Ancestors(parentBranchName)
	return &appendConfig{
		ancestorBranches:    ancestorBranches,
		isOffline:           isOffline,
		hasOrigin:           hasOrigin,
		mainBranch:          mainBranch,
		noPushHook:          !pushHook,
		parentBranch:        *parentBranch,
		shouldNewBranchPush: shouldNewBranchPush,
		targetBranch:        targetBranch,
	}, fc.Err
}

func appendStepList(config *appendConfig, run *git.ProdRunner) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range append(config.ancestorBranches, config.parentBranch) {
		updateBranchSteps(&list, branch.Name, true, run)
	}
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch.Name})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch.Name})
	list.Add(&steps.CheckoutStep{Branch: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		list.Add(&steps.CreateTrackingBranchStep{Branch: config.targetBranch, NoPushHook: config.noPushHook})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, &run.Backend, config.mainBranch)
	return list.Result()
}
