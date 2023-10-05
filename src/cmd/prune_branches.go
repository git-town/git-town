package cmd

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/runvm"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/spf13/cobra"
)

const pruneBranchesDesc = "Deletes local branches whose tracking branch no longer exists"

const pruneBranchesHelp = `
Deletes branches whose tracking branch no longer exists from the local repository.
This usually means the branch was shipped or killed on another machine.`

func pruneBranchesCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "prune-branches",
		Args:  cobra.NoArgs,
		Short: pruneBranchesDesc,
		Long:  long(pruneBranchesDesc, pruneBranchesHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePruneBranches(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func executePruneBranches(debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: true,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determinePruneBranchesConfig(repo, debug)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "prune-branches",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunSteps:            pruneBranchesSteps(config),
	}
	return runvm.Execute(runvm.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               nil,
		Debug:                   debug,
		Lineage:                 config.lineage,
		NoPushHook:              !config.pushHook,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
	})
}

type pruneBranchesConfig struct {
	branches                  domain.Branches
	lineage                   config.Lineage
	branchesWithDeletedRemote domain.LocalBranchNames
	mainBranch                domain.LocalBranchName
	previousBranch            domain.LocalBranchName
	pushHook                  bool
}

func determinePruneBranchesConfig(repo *execute.OpenRepoResult, debug bool) (*pruneBranchesConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage()
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Debug:                 debug,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	return &pruneBranchesConfig{
		branches:                  branches,
		lineage:                   lineage,
		branchesWithDeletedRemote: branches.All.LocalBranchesWithDeletedTrackingBranches().Names(),
		mainBranch:                repo.Runner.Config.MainBranch(),
		previousBranch:            repo.Runner.Backend.PreviouslyCheckedOutBranch(),
		pushHook:                  pushHook,
	}, branchesSnapshot, stashSnapshot, exit, err
}

func pruneBranchesSteps(config *pruneBranchesConfig) steps.List {
	result := steps.List{}
	for _, branchWithDeletedRemote := range config.branchesWithDeletedRemote {
		if config.branches.Initial == branchWithDeletedRemote {
			result.Add(&step.Checkout{Branch: config.mainBranch})
		}
		parent := config.lineage.Parent(branchWithDeletedRemote)
		for _, child := range config.lineage.Children(branchWithDeletedRemote) {
			if parent.IsEmpty() {
				result.Add(&step.DeleteParentBranch{Branch: child})
			} else {
				result.Add(&step.SetParent{Branch: child, ParentBranch: parent})
			}
		}
		if config.branches.Types.IsFeatureBranch(branchWithDeletedRemote) {
			result.Add(&step.DeleteParentBranch{Branch: branchWithDeletedRemote})
		}
		if config.branches.Types.IsPerennialBranch(branchWithDeletedRemote) {
			result.Add(&step.RemoveFromPerennialBranches{Branch: branchWithDeletedRemote})
		}
		result.Add(&step.DeleteLocalBranch{Branch: branchWithDeletedRemote, Parent: config.mainBranch.Location(), Force: false})
	}
	result.Wrap(steps.WrapOptions{
		RunInGitRoot:     false,
		StashOpenChanges: false,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return result
}
