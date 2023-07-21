package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/validate"
	"github.com/spf13/cobra"
)

const killDesc = "Removes an obsolete feature branch"

const killHelp = `
Deletes the current or provided branch from the local and origin repositories.
Does not delete perennial branches nor the main branch.`

func killCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "kill [<branch>]",
		Args:  cobra.MaximumNArgs(1),
		Short: killDesc,
		Long:  long(killDesc, killHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return kill(args, readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func kill(args []string, debug bool) error {
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: false,
		ValidateGitversion:    true,
		ValidateIsConfigured:  true,
		ValidateIsOnline:      false,
		ValidateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineKillConfig(args, &run)
	if err != nil {
		return err
	}
	stepList, err := killStepList(config, &run)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "kill",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &run, nil)
}

type killConfig struct {
	childBranches       []string
	hasOpenChanges      bool
	hasTrackingBranch   bool
	initialBranch       string
	isOffline           bool
	isTargetBranchLocal bool
	mainBranch          string
	noPushHook          bool
	previousBranch      string
	targetBranchParent  string
	targetBranch        string
}

func determineKillConfig(args []string, run *git.ProdRunner) (*killConfig, error) {
	initialBranch, err := run.Backend.CurrentBranch()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	var targetBranch string
	if len(args) > 0 {
		targetBranch = args[0]
	} else {
		targetBranch = initialBranch
	}
	if !run.Config.IsFeatureBranch(targetBranch) {
		return nil, fmt.Errorf("you can only kill feature branches")
	}
	isTargetBranchLocal, err := run.Backend.HasLocalBranch(targetBranch)
	if err != nil {
		return nil, err
	}
	if isTargetBranchLocal {
		err = validate.KnowsBranchAncestors(targetBranch, mainBranch, &run.Backend)
		if err != nil {
			return nil, err
		}
		run.Config.Reload()
	}
	hasOrigin, err := run.Backend.HasOrigin()
	if err != nil {
		return nil, err
	}
	isOffline, err := run.Config.IsOffline()
	if err != nil {
		return nil, err
	}
	if hasOrigin && !isOffline {
		err := run.Frontend.Fetch()
		if err != nil {
			return nil, err
		}
	}
	if initialBranch != targetBranch {
		hasTargetBranch, err := run.Backend.HasLocalOrOriginBranch(targetBranch, mainBranch)
		if err != nil {
			return nil, err
		}
		if !hasTargetBranch {
			return nil, fmt.Errorf("there is no branch named %q", targetBranch)
		}
	}
	hasTrackingBranch, err := run.Backend.HasTrackingBranch(targetBranch)
	if err != nil {
		return nil, err
	}
	previousBranch, err := run.Backend.PreviouslyCheckedOutBranch()
	if err != nil {
		return nil, err
	}
	hasOpenChanges, err := run.Backend.HasOpenChanges()
	if err != nil {
		return nil, err
	}
	pushHook, err := run.Config.PushHook()
	if err != nil {
		return nil, err
	}
	lineage := run.Config.Lineage()
	return &killConfig{
		childBranches:       lineage.Children(targetBranch),
		hasOpenChanges:      hasOpenChanges,
		hasTrackingBranch:   hasTrackingBranch,
		initialBranch:       initialBranch,
		isOffline:           isOffline,
		isTargetBranchLocal: isTargetBranchLocal,
		mainBranch:          mainBranch,
		noPushHook:          !pushHook,
		previousBranch:      previousBranch,
		targetBranch:        targetBranch,
		targetBranchParent:  lineage.Parent(targetBranch),
	}, nil
}

func killStepList(config *killConfig, run *git.ProdRunner) (runstate.StepList, error) {
	result := runstate.StepList{}
	switch {
	case config.isTargetBranchLocal:
		if config.hasTrackingBranch && !config.isOffline {
			result.Append(&steps.DeleteOriginBranchStep{Branch: config.targetBranch, IsTracking: true, NoPushHook: config.noPushHook})
		}
		if config.initialBranch == config.targetBranch {
			if config.hasOpenChanges {
				result.Append(&steps.CommitOpenChangesStep{})
			}
			result.Append(&steps.CheckoutStep{Branch: config.targetBranchParent})
		}
		result.Append(&steps.DeleteLocalBranchStep{Branch: config.targetBranch, Parent: config.mainBranch, Force: true})
		for _, child := range config.childBranches {
			result.Append(&steps.SetParentStep{Branch: child, ParentBranch: config.targetBranchParent})
		}
		result.Append(&steps.DeleteParentBranchStep{Branch: config.targetBranch, Parent: run.Config.Lineage().Parent(config.targetBranch)})
	case !config.isOffline:
		result.Append(&steps.DeleteOriginBranchStep{Branch: config.targetBranch, IsTracking: false, NoPushHook: config.noPushHook})
	default:
		return runstate.StepList{}, fmt.Errorf("cannot delete remote branch %q in offline mode", config.targetBranch)
	}
	err := result.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.initialBranch != config.targetBranch && config.targetBranch == config.previousBranch,
	}, &run.Backend, config.mainBranch)
	return result, err
}
