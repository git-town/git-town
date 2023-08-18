package cmd

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/runstate"
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
			return pruneBranches(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func pruneBranches(debug bool) error {
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 true,
		HandleUnfinishedState: true,
		OmitBranchNames:       false,
		ValidateIsOnline:      true,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	config, err := determinePruneBranchesConfig(&repo.Runner)
	if err != nil {
		return err
	}
	stepList, err := pruneBranchesStepList(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "prune-branches",
		RunStepList: stepList,
	}
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &runState,
		Run:       &repo.Runner,
		Connector: nil,
		RootDir:   repo.RootDir,
	})
}

type pruneBranchesConfig struct {
	branchDurations  config.BranchDurations
	initialBranch    domain.LocalBranchName
	lineage          config.Lineage
	branchesToDelete domain.LocalBranchNames
	mainBranch       domain.LocalBranchName
	previousBranch   domain.LocalBranchName
}

func determinePruneBranchesConfig(run *git.ProdRunner) (*pruneBranchesConfig, error) {
	branches, err := execute.LoadBranches(run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	return &pruneBranchesConfig{
		branchDurations:  branches.Durations,
		initialBranch:    branches.Initial,
		lineage:          run.Config.Lineage(),
		branchesToDelete: branches.All.LocalBranchesWithDeletedTrackingBranches().Names(),
		mainBranch:       run.Config.MainBranch(),
		previousBranch:   run.Backend.PreviouslyCheckedOutBranch(),
	}, err
}

func pruneBranchesStepList(config *pruneBranchesConfig) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branchWithDeletedRemote := range config.branchesToDelete {
		if config.initialBranch == branchWithDeletedRemote {
			result.Append(&steps.CheckoutStep{Branch: config.mainBranch})
		}
		parent := config.lineage.Parent(branchWithDeletedRemote)
		if !parent.IsEmpty() {
			for _, child := range config.lineage.Children(branchWithDeletedRemote) {
				result.Append(&steps.SetParentStep{Branch: child, ParentBranch: parent})
			}
			result.Append(&steps.DeleteParentBranchStep{Branch: branchWithDeletedRemote, Parent: config.lineage.Parent(branchWithDeletedRemote)})
		}
		if config.branchDurations.IsPerennialBranch(branchWithDeletedRemote) {
			result.Append(&steps.RemoveFromPerennialBranchesStep{Branch: branchWithDeletedRemote})
		}
		result.Append(&steps.DeleteLocalBranchStep{Branch: branchWithDeletedRemote, Parent: config.mainBranch.Location})
	}
	err := result.Wrap(runstate.WrapOptions{
		RunInGitRoot:     false,
		StashOpenChanges: false,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranch,
		PreviousBranch:   config.previousBranch,
	})
	return result, err
}
