package cmd

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
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
	repo, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: true,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, exit, err := determinePruneBranchesConfig(&repo)
	if err != nil || exit {
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
		Branches:  config.branches.All,
	})
}

type pruneBranchesConfig struct {
	branches         domain.Branches
	lineage          config.Lineage
	branchesToDelete domain.LocalBranchNames
	mainBranch       domain.LocalBranchName
	previousBranch   domain.LocalBranchName
}

func determinePruneBranchesConfig(repo *execute.RepoData) (*pruneBranchesConfig, bool, error) {
	branches, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	return &pruneBranchesConfig{
		branches:         branches,
		lineage:          repo.Runner.Config.Lineage(),
		branchesToDelete: branches.All.LocalBranchesWithDeletedTrackingBranches().Names(),
		mainBranch:       repo.Runner.Config.MainBranch(),
		previousBranch:   repo.Runner.Backend.PreviouslyCheckedOutBranch(),
	}, exit, err
}

func pruneBranchesStepList(config *pruneBranchesConfig) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branchWithDeletedRemote := range config.branchesToDelete {
		if config.branches.Initial == branchWithDeletedRemote {
			result.Append(&steps.CheckoutStep{Branch: config.mainBranch})
		}
		parent := config.lineage.Parent(branchWithDeletedRemote)
		if !parent.IsEmpty() {
			for _, child := range config.lineage.Children(branchWithDeletedRemote) {
				result.Append(&steps.SetParentStep{Branch: child, ParentBranch: parent})
			}
			result.Append(&steps.DeleteParentBranchStep{Branch: branchWithDeletedRemote, Parent: config.lineage.Parent(branchWithDeletedRemote)})
		}
		if config.branches.Types.IsPerennialBranch(branchWithDeletedRemote) {
			result.Append(&steps.RemoveFromPerennialBranchesStep{Branch: branchWithDeletedRemote})
		}
		result.Append(&steps.DeleteLocalBranchStep{Branch: branchWithDeletedRemote, Parent: config.mainBranch.Location(), Force: false})
	}
	err := result.Wrap(runstate.WrapOptions{
		RunInGitRoot:     false,
		StashOpenChanges: false,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return result, err
}
