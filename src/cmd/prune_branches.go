package cmd

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/spf13/cobra"
)

const pruneBranchesDesc = "Deletes local branches whose tracking branch no longer exists"

const pruneBranchesHelp = `
Deletes branches whose tracking branch no longer exists from the local repository.
This usually means the branch was shipped or killed on another machine.`

func pruneBranchesCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := debugFlag()
	cmd := cobra.Command{
		Use:   "prune-branches",
		Args:  cobra.NoArgs,
		Short: pruneBranchesDesc,
		Long:  long(pruneBranchesDesc, pruneBranchesHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPruneBranches(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func runPruneBranches(debug bool) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		debug:                 debug,
		dryRun:                false,
		handleUnfinishedState: true,
		validateGitversion:    true,
		validateIsRepository:  true,
		validateIsConfigured:  true,
		validateIsOnline:      true,
	})
	if err != nil || exit {
		return err
	}
	config, err := determinePruneBranchesConfig(&repo)
	if err != nil {
		return err
	}
	stepList, err := pruneBranchesStepList(config, &repo)
	if err != nil {
		return err
	}
	runState := runstate.New("prune-branches", stepList)
	return runstate.Execute(runState, &repo, nil)
}

type pruneBranchesConfig struct {
	initialBranch                            string
	localBranchesWithDeletedTrackingBranches []string
	mainBranch                               string
}

func determinePruneBranchesConfig(repo *git.PublicRepo) (*pruneBranchesConfig, error) {
	hasOrigin, err := repo.HasOrigin()
	if err != nil {
		return nil, err
	}
	if hasOrigin {
		err = repo.Fetch()
		if err != nil {
			return nil, err
		}
	}
	initialBranch, err := repo.CurrentBranch()
	if err != nil {
		return nil, err
	}
	localBranchesWithDeletedTrackingBranches, err := repo.LocalBranchesWithDeletedTrackingBranches()
	if err != nil {
		return nil, err
	}
	return &pruneBranchesConfig{
		initialBranch:                            initialBranch,
		localBranchesWithDeletedTrackingBranches: localBranchesWithDeletedTrackingBranches,
		mainBranch:                               repo.Config.MainBranch(),
	}, nil
}

func pruneBranchesStepList(config *pruneBranchesConfig, repo *git.PublicRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branchWithDeletedRemote := range config.localBranchesWithDeletedTrackingBranches {
		if config.initialBranch == branchWithDeletedRemote {
			result.Append(&steps.CheckoutStep{Branch: config.mainBranch})
		}
		parent := repo.Config.ParentBranch(branchWithDeletedRemote)
		if parent != "" {
			for _, child := range repo.Config.ChildBranches(branchWithDeletedRemote) {
				result.Append(&steps.SetParentStep{Branch: child, ParentBranch: parent})
			}
			result.Append(&steps.DeleteParentBranchStep{Branch: branchWithDeletedRemote})
		}
		if repo.Config.IsPerennialBranch(branchWithDeletedRemote) {
			result.Append(&steps.RemoveFromPerennialBranchesStep{Branch: branchWithDeletedRemote})
		}
		result.Append(&steps.DeleteLocalBranchStep{Branch: branchWithDeletedRemote, Parent: config.mainBranch})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: false, StashOpenChanges: false}, &repo.InternalRepo, config.mainBranch)
	return result, err
}
