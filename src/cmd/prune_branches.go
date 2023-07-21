package cmd

import (
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
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: false,
	})
	if err != nil {
		return err
	}
	_, initialBranch, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		HandleUnfinishedState: true,
		ValidateGitversion:    true,
		ValidateIsRepository:  true,
		ValidateIsConfigured:  true,
		ValidateIsOnline:      true,
	})
	if err != nil || exit {
		return err
	}
	config, err := determinePruneBranchesConfig(&run, initialBranch)
	if err != nil {
		return err
	}
	stepList, err := pruneBranchesStepList(config, &run)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "prune-branches",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &run, nil)
}

type pruneBranchesConfig struct {
	initialBranch                            string
	localBranchesWithDeletedTrackingBranches []string
	mainBranch                               string
}

func determinePruneBranchesConfig(run *git.ProdRunner, initialBranch string) (*pruneBranchesConfig, error) {
	hasOrigin, err := run.Backend.HasOrigin()
	if err != nil {
		return nil, err
	}
	if hasOrigin {
		err = run.Frontend.Fetch()
		if err != nil {
			return nil, err
		}
	}
	localBranchesWithDeletedTrackingBranches, err := run.Backend.LocalBranchesWithDeletedTrackingBranches()
	if err != nil {
		return nil, err
	}
	return &pruneBranchesConfig{
		initialBranch:                            initialBranch,
		localBranchesWithDeletedTrackingBranches: localBranchesWithDeletedTrackingBranches,
		mainBranch:                               run.Config.MainBranch(),
	}, nil
}

func pruneBranchesStepList(config *pruneBranchesConfig, run *git.ProdRunner) (runstate.StepList, error) {
	result := runstate.StepList{}
	lineage := run.Config.Lineage()
	for _, branchWithDeletedRemote := range config.localBranchesWithDeletedTrackingBranches {
		if config.initialBranch == branchWithDeletedRemote {
			result.Append(&steps.CheckoutStep{Branch: config.mainBranch})
		}
		parent := lineage.Parent(branchWithDeletedRemote)
		if parent != "" {
			for _, child := range lineage.Children(branchWithDeletedRemote) {
				result.Append(&steps.SetParentStep{Branch: child, ParentBranch: parent})
			}
			result.Append(&steps.DeleteParentBranchStep{Branch: branchWithDeletedRemote, Parent: lineage.Parent(branchWithDeletedRemote)})
		}
		if run.Config.IsPerennialBranch(branchWithDeletedRemote) {
			result.Append(&steps.RemoveFromPerennialBranchesStep{Branch: branchWithDeletedRemote})
		}
		result.Append(&steps.DeleteLocalBranchStep{Branch: branchWithDeletedRemote, Parent: config.mainBranch})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: false, StashOpenChanges: false}, &run.Backend, config.mainBranch)
	return result, err
}
