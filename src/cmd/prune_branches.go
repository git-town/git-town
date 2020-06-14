package cmd

import (
	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/steps"
	"github.com/spf13/cobra"
)

type pruneBranchesConfig struct {
	initialBranchName                        string
	mainBranch                               string
	localBranchesWithDeletedTrackingBranches []string
}

var pruneBranchesCommand = &cobra.Command{
	Use:   "prune-branches",
	Short: "Deletes local branches whose tracking branch no longer exists",
	Long: `Deletes local branches whose tracking branch no longer exists

Deletes branches whose tracking branch no longer exists from the local repository.
This usually means the branch was shipped or killed on another machine.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getPruneBranchesConfig(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		stepList, err := getPruneBranchesStepList(config, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		runState := steps.NewRunState("prune-branches", stepList)
		err = steps.Run(runState, prodRepo, nil)
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
		return git.Config().ValidateIsOnline()
	},
}

func getPruneBranchesConfig(repo *git.ProdRepo) (result pruneBranchesConfig, err error) {
	hasOrigin, err := repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	if hasOrigin {
		err = repo.Logging.Fetch()
		if err != nil {
			return result, err
		}
	}
	result.mainBranch = git.Config().GetMainBranch()
	result.initialBranchName, err = repo.Silent.CurrentBranch()
	if err != nil {
		return result, err
	}
	result.localBranchesWithDeletedTrackingBranches, err = repo.Silent.LocalBranchesWithDeletedTrackingBranches()
	return result, err
}

func getPruneBranchesStepList(config pruneBranchesConfig, repo *git.ProdRepo) (result steps.StepList, err error) {
	initialBranchName := config.initialBranchName
	for _, branchName := range config.localBranchesWithDeletedTrackingBranches {
		if initialBranchName == branchName {
			result.Append(&steps.CheckoutBranchStep{BranchName: config.mainBranch})
		}
		parent := git.Config().GetParentBranch(branchName)
		if parent != "" {
			for _, child := range git.Config().GetChildBranches(branchName) {
				result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: parent})
			}
			result.Append(&steps.DeleteParentBranchStep{BranchName: branchName})
		}
		if git.Config().IsPerennialBranch(branchName) {
			result.Append(&steps.RemoveFromPerennialBranches{BranchName: branchName})
		}
		result.Append(&steps.DeleteLocalBranchStep{BranchName: branchName})
	}
	err = result.Wrap(steps.WrapOptions{RunInGitRoot: false, StashOpenChanges: false}, repo)
	return result, err
}

func init() {
	RootCmd.AddCommand(pruneBranchesCommand)
}
