package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
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
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: false,
	})
	if err != nil {
		return err
	}
	allBranches, initialBranch, rootDir, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 true,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  true,
		ValidateIsOnline:      false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineKillConfig(args, &run, allBranches, initialBranch)
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
	return runstate.Execute(&runState, &run, nil, rootDir)
}

type killConfig struct {
	childBranches      []string
	hasOpenChanges     bool
	initialBranch      string
	isOffline          bool
	mainBranch         string
	noPushHook         bool
	previousBranch     string
	targetBranchParent string
	targetBranch       git.BranchSyncStatus
}

func determineKillConfig(args []string, run *git.ProdRunner, allBranches git.BranchesSyncStatus, initialBranch string) (*killConfig, error) {
	mainBranch := run.Config.MainBranch()
	targetBranchName := initialBranch
	if len(args) > 0 {
		targetBranchName = args[0]
	}
	if !run.Config.IsFeatureBranch(targetBranchName) {
		return nil, fmt.Errorf(messages.KillOnlyFeatureBranches)
	}
	targetBranch := allBranches.Lookup(targetBranchName)
	if targetBranch == nil {
		return nil, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	if targetBranch.IsLocal() {
		err := validate.KnowsBranchAncestors(targetBranchName, mainBranch, &run.Backend)
		if err != nil {
			return nil, err
		}
		run.Config.Reload()
	}
	isOffline, err := run.Config.IsOffline()
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
		childBranches:      lineage.Children(targetBranchName),
		hasOpenChanges:     hasOpenChanges,
		initialBranch:      initialBranch,
		isOffline:          isOffline,
		mainBranch:         mainBranch,
		noPushHook:         !pushHook,
		previousBranch:     previousBranch,
		targetBranch:       *targetBranch,
		targetBranchParent: lineage.Parent(targetBranchName),
	}, nil
}

func killStepList(config *killConfig, run *git.ProdRunner) (runstate.StepList, error) {
	result := runstate.StepList{}
	switch {
	case config.targetBranch.IsLocal():
		if config.targetBranch.HasTrackingBranch() && !config.isOffline {
			result.Append(&steps.DeleteOriginBranchStep{Branch: config.targetBranch.Name, IsTracking: true, NoPushHook: config.noPushHook})
		}
		if config.initialBranch == config.targetBranch.Name {
			if config.hasOpenChanges {
				result.Append(&steps.CommitOpenChangesStep{})
			}
			result.Append(&steps.CheckoutStep{Branch: config.targetBranchParent})
		}
		result.Append(&steps.DeleteLocalBranchStep{Branch: config.targetBranch.Name, Parent: config.mainBranch, Force: true})
		for _, child := range config.childBranches {
			result.Append(&steps.SetParentStep{Branch: child, ParentBranch: config.targetBranchParent})
		}
		result.Append(&steps.DeleteParentBranchStep{Branch: config.targetBranch.Name, Parent: config.targetBranchParent})
	case !config.isOffline:
		result.Append(&steps.DeleteOriginBranchStep{Branch: config.targetBranch.Name, IsTracking: false, NoPushHook: config.noPushHook})
	default:
		return runstate.StepList{}, fmt.Errorf(messages.DeleteRemoteBranchOffline, config.targetBranch.Name)
	}
	err := result.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.initialBranch != config.targetBranch.Name && config.targetBranch.Name == config.previousBranch,
	}, &run.Backend, config.mainBranch)
	return result, err
}
