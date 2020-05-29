package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
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
		config, err := getPruneBranchesConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList := getPruneBranchesStepList(config)
		runState := steps.NewRunState("prune-branches", stepList)
		err = steps.Run(runState)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := git.ValidateIsRepository(); err != nil {
			return err
		}
		if err := validateIsConfigured(); err != nil {
			return err
		}
		return git.Config().ValidateIsOnline()
	},
}

func getPruneBranchesConfig() (result pruneBranchesConfig, err error) {
	if git.HasRemote("origin") {
		err = script.Fetch()
		if err != nil {
			return result, err
		}
	}
	result.mainBranch = git.Config().GetMainBranch()
	result.initialBranchName = git.GetCurrentBranchName()
	result.localBranchesWithDeletedTrackingBranches = git.GetLocalBranchesWithDeletedTrackingBranches()
	return result, nil
}

func getPruneBranchesStepList(config pruneBranchesConfig) (result steps.StepList) {
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
	result.Wrap(steps.WrapOptions{RunInGitRoot: false, StashOpenChanges: false})
	return
}

func init() {
	RootCmd.AddCommand(pruneBranchesCommand)
}
