package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/slice"
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
	config, exit, err := determineKillConfig(args, &repo)
	if err != nil || exit {
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
	initialBranch  domain.LocalBranchName
	isOffline      bool
	lineage        config.Lineage
	mainBranch     domain.LocalBranchName
	noPushHook     bool
	previousBranch domain.LocalBranchName
	targetBranch   domain.BranchInfo
}

func determineKillConfig(args []string, repo *execute.OpenRepoResult) (*killConfig, bool, error) {
	branches, exit, err := execute.LoadSnapshot(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 true,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	targetBranchName := domain.NewLocalBranchName(slice.FirstElementOr(args, branches.Initial.String()))
	if !branches.Types.IsFeatureBranch(targetBranchName) {
		return nil, false, fmt.Errorf(messages.KillOnlyFeatureBranches)
	}
	targetBranch := branches.All.FindLocalBranch(targetBranchName)
	if targetBranch == nil {
		return nil, false, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	lineage := repo.Runner.Config.Lineage()
	if targetBranch.IsLocal() {
		updated, err := validate.KnowsBranchAncestors(targetBranchName, validate.KnowsBranchAncestorsArgs{
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
			repo.Runner.Config.Reload()
			lineage = repo.Runner.Config.Lineage()
		}
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges, err := repo.Runner.Backend.HasOpenChanges()
	if err != nil {
		return nil, false, err
	}
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, false, err
	}
	return &killConfig{
		hasOpenChanges: hasOpenChanges,
		initialBranch:  branches.Initial,
		isOffline:      repo.IsOffline,
		lineage:        lineage,
		mainBranch:     mainBranch,
		noPushHook:     !pushHook,
		previousBranch: previousBranch,
		targetBranch:   *targetBranch,
	}, false, nil
}

func (kc killConfig) isOnline() bool {
	return !kc.isOffline
}

func (kc killConfig) targetBranchParent() domain.LocalBranchName {
	return kc.lineage.Parent(kc.targetBranch.LocalName)
}

func killStepList(config *killConfig) (runstate.StepList, error) {
	result := runstate.StepList{}
	killFeatureBranch(&result, *config)
	err := result.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.initialBranch != config.targetBranch.LocalName && config.targetBranch.LocalName == config.previousBranch && config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranch,
		PreviousBranch:   config.previousBranch,
	})
	return result, err
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killFeatureBranch(list *runstate.StepList, config killConfig) {
	if config.targetBranch.HasTrackingBranch() && config.isOnline() {
		list.Append(&steps.DeleteTrackingBranchStep{Branch: config.targetBranch.LocalName, NoPushHook: config.noPushHook})
	}
	if config.initialBranch == config.targetBranch.LocalName {
		if config.hasOpenChanges {
			list.Append(&steps.CommitOpenChangesStep{})
		}
		list.Append(&steps.CheckoutStep{Branch: config.targetBranchParent()})
	}
	list.Append(&steps.DeleteLocalBranchStep{Branch: config.targetBranch.LocalName, Parent: config.mainBranch.Location(), Force: true})
	childBranches := config.lineage.Children(config.targetBranch.LocalName)
	for _, child := range childBranches {
		list.Append(&steps.SetParentStep{Branch: child, ParentBranch: config.targetBranchParent()})
	}
	list.Append(&steps.DeleteParentBranchStep{Branch: config.targetBranch.LocalName, Parent: config.targetBranchParent()})
}
