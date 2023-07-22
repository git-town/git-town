package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
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
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: false,
	})
	if err != nil {
		return err
	}
	branchesSyncStatus, initialBranch, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateIsOnline:      false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	config, err := determinePrependConfig(args, &run, branchesSyncStatus, initialBranch)
	if err != nil {
		return err
	}
	stepList, err := prependStepList(config, &run)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "prepend",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &run, nil)
}

type prependConfig struct {
	ancestorBranches    []string
	hasOrigin           bool
	initialBranch       string
	isOffline           bool
	mainBranch          string
	noPushHook          bool
	parentBranch        string
	shouldNewBranchPush bool
	targetBranch        string
}

func determinePrependConfig(args []string, run *git.ProdRunner, branchesSyncStatus git.BranchesSyncStatus, initialBranch string) (*prependConfig, error) {
	fc := failure.Collector{}
	hasOrigin := fc.Bool(run.Backend.HasOrigin())
	shouldNewBranchPush := fc.Bool(run.Config.ShouldNewBranchPush())
	pushHook := fc.Bool(run.Config.PushHook())
	isOffline := fc.Bool(run.Config.IsOffline())
	mainBranch := run.Config.MainBranch()
	if fc.Err != nil {
		return nil, fc.Err
	}
	targetBranch := args[0]
	if branchesSyncStatus.Contains(targetBranch) {
		return nil, fmt.Errorf("a branch named %q already exists", targetBranch)
	}
	if !run.Config.IsFeatureBranch(initialBranch) {
		return nil, fmt.Errorf("the branch %q is not a feature branch. Only feature branches can have parent branches", initialBranch)
	}
	err := validate.KnowsBranchAncestors(initialBranch, run.Config.MainBranch(), &run.Backend)
	if err != nil {
		return nil, err
	}
	lineage := run.Config.Lineage()
	return &prependConfig{
		hasOrigin:           hasOrigin,
		initialBranch:       initialBranch,
		isOffline:           isOffline,
		mainBranch:          mainBranch,
		noPushHook:          !pushHook,
		parentBranch:        lineage.Parent(initialBranch),
		ancestorBranches:    lineage.Ancestors(initialBranch),
		shouldNewBranchPush: shouldNewBranchPush,
		targetBranch:        targetBranch,
	}, nil
}

func prependStepList(config *prependConfig, run *git.ProdRunner) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.ancestorBranches {
		updateBranchSteps(&list, branch, true, run)
	}
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.initialBranch, ParentBranch: config.targetBranch})
	list.Add(&steps.CheckoutStep{Branch: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		list.Add(&steps.CreateTrackingBranchStep{Branch: config.targetBranch, NoPushHook: config.noPushHook})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, &run.Backend, config.mainBranch)
	return list.Result()
}
