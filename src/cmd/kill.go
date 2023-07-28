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
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 true,
		HandleUnfinishedState: false,
		OmitBranchNames:       false,
		ValidateIsOnline:      false,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	branches, err := execute.LoadBranches(&repo.Runner, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	config, err := determineKillConfig(args, &repo.Runner, branches, repo.IsOffline)
	if err != nil {
		return err
	}
	stepList, err := killStepList(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "kill",
		RunStepList: stepList,
	}
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &runState,
		Run:       &repo.Runner,
		Connector: nil,
		RootDir:   repo.RootDir,
	})
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

func determineKillConfig(args []string, run *git.ProdRunner, branches execute.Branches, isOffline bool) (*killConfig, error) {
	mainBranch := run.Config.MainBranch()
	targetBranchName := branches.Initial
	if len(args) > 0 {
		targetBranchName = args[0]
	}
	if !branches.Durations.IsFeatureBranch(targetBranchName) {
		return nil, fmt.Errorf(messages.KillOnlyFeatureBranches)
	}
	targetBranch := branches.All.Lookup(targetBranchName)
	if targetBranch == nil {
		return nil, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	lineage := run.Config.Lineage()
	if targetBranch.IsLocal() {
		updated, err := validate.KnowsBranchAncestors(targetBranchName, validate.KnowsBranchAncestorsArgs{
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
			run.Config.Reload()
			lineage = run.Config.Lineage()
		}
	}
	previousBranch := run.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges, err := run.Backend.HasOpenChanges()
	if err != nil {
		return nil, err
	}
	pushHook, err := run.Config.PushHook()
	if err != nil {
		return nil, err
	}
	return &killConfig{
		childBranches:      lineage.Children(targetBranchName),
		hasOpenChanges:     hasOpenChanges,
		initialBranch:      branches.Initial,
		isOffline:          isOffline,
		mainBranch:         mainBranch,
		noPushHook:         !pushHook,
		previousBranch:     previousBranch,
		targetBranch:       *targetBranch,
		targetBranchParent: lineage.Parent(targetBranchName),
	}, nil
}

func killStepList(config *killConfig) (runstate.StepList, error) {
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
		StashOpenChanges: config.initialBranch != config.targetBranch.Name && config.targetBranch.Name == config.previousBranch && config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranch,
		PreviousBranch:   config.previousBranch,
	})
	return result, err
}
