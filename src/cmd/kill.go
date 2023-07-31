package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
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
	config, err := determineKillConfig(args, &repo.Runner, repo.IsOffline)
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
	hasOpenChanges bool
	initialBranch  string
	isOffline      bool
	lineage        config.Lineage
	mainBranch     string
	noPushHook     bool
	previousBranch string
	targetBranch   git.BranchSyncStatus
}

func determineKillConfig(args []string, run *git.ProdRunner, isOffline bool) (*killConfig, error) {
	branches, err := execute.LoadBranches(run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return nil, err
	}
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
		hasOpenChanges: hasOpenChanges,
		initialBranch:  branches.Initial,
		isOffline:      isOffline,
		lineage:        lineage,
		mainBranch:     mainBranch,
		noPushHook:     !pushHook,
		previousBranch: previousBranch,
		targetBranch:   *targetBranch,
	}, nil
}

func (kc killConfig) isOnline() bool {
	return !kc.isOffline
}

func (kc killConfig) targetBranchParent() string {
	return kc.lineage.Parent(kc.targetBranch.Name)
}

func killStepList(config *killConfig) (runstate.StepList, error) {
	result := runstate.StepList{}
	if config.targetBranch.IsLocal() {
		killFeatureBranch(&result, *config)
	} else if config.isOnline() {
		// user wants us to kill a remote branch and we are online
		result.Append(&steps.DeleteOriginBranchStep{Branch: config.targetBranch.Name, IsTracking: false, NoPushHook: config.noPushHook})
	} else {
		// user wants us to kill a remote branch and we are offline
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

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely)
func killFeatureBranch(list *runstate.StepList, config killConfig) {
	if config.targetBranch.HasTrackingBranch() && config.isOnline() {
		list.Append(&steps.DeleteOriginBranchStep{Branch: config.targetBranch.Name, IsTracking: true, NoPushHook: config.noPushHook})
	}
	if config.initialBranch == config.targetBranch.Name {
		if config.hasOpenChanges {
			list.Append(&steps.CommitOpenChangesStep{})
		}
		list.Append(&steps.CheckoutStep{Branch: config.targetBranchParent()})
	}
	list.Append(&steps.DeleteLocalBranchStep{Branch: config.targetBranch.Name, Parent: config.mainBranch, Force: true})
	childBranches := config.lineage.Children(config.targetBranch.Name)
	for _, child := range childBranches {
		list.Append(&steps.SetParentStep{Branch: child, ParentBranch: config.targetBranchParent()})
	}
	list.Append(&steps.DeleteParentBranchStep{Branch: config.targetBranch.Name, Parent: config.targetBranchParent()})
}
