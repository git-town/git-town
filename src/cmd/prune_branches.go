package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/spf13/cobra"
)

type pruneBranchesConfig struct {
	initialBranchName                        string
	localBranchesWithDeletedTrackingBranches []string
	mainBranch                               string
}

var pruneBranchesCommand = &cobra.Command{
	Use:   "prune-branches",
	Short: "Deletes local branches whose tracking branch no longer exists",
	Long: `Deletes local branches whose tracking branch no longer exists

Deletes branches whose tracking branch no longer exists from the local repository.
This usually means the branch was shipped or killed on another machine.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := createPruneBranchesConfig(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		stepList, err := createPruneBranchesStepList(config, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		runState := runstate.New("prune-branches", stepList)
		err = runstate.Execute(runState, prodRepo, nil)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateIsRepository(prodRepo); err != nil {
			return err
		}
		if err := validateIsConfigured(prodRepo); err != nil {
			return err
		}
		return prodRepo.Config.ValidateIsOnline()
	},
}

func createPruneBranchesConfig(repo *git.ProdRepo) (pruneBranchesConfig, error) {
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return pruneBranchesConfig{}, err
	}
	if hasOrigin {
		err = repo.Logging.Fetch()
		if err != nil {
			return pruneBranchesConfig{}, err
		}
	}
	result := pruneBranchesConfig{
		mainBranch: repo.Config.MainBranch(),
	}
	result.initialBranchName, err = repo.Silent.CurrentBranch()
	if err != nil {
		return pruneBranchesConfig{}, err
	}
	result.localBranchesWithDeletedTrackingBranches, err = repo.Silent.LocalBranchesWithDeletedTrackingBranches()
	return result, err
}

func createPruneBranchesStepList(config pruneBranchesConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	initialBranchName := config.initialBranchName
	result := runstate.StepList{}
	for _, branchName := range config.localBranchesWithDeletedTrackingBranches {
		if initialBranchName == branchName {
			result.Append(&steps.CheckoutBranchStep{BranchName: config.mainBranch})
		}
		parent := repo.Config.ParentBranch(branchName)
		if parent != "" {
			for _, child := range repo.Config.ChildBranches(branchName) {
				result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: parent})
			}
			result.Append(&steps.DeleteParentBranchStep{BranchName: branchName})
		}
		if repo.Config.IsPerennialBranch(branchName) {
			result.Append(&steps.RemoveFromPerennialBranchesStep{BranchName: branchName})
		}
		result.Append(&steps.DeleteLocalBranchStep{BranchName: branchName})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: false, StashOpenChanges: false}, repo)
	return result, err
}

func init() {
	RootCmd.AddCommand(pruneBranchesCommand)
}
