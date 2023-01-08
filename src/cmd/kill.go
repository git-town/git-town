package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/spf13/cobra"
)

type killConfig struct {
	childBranches       []string
	hasOpenChanges      bool
	hasTrackingBranch   bool
	initialBranch       string
	isOffline           bool
	isTargetBranchLocal bool
	noPushVerify        bool
	previousBranch      string
	targetBranchParent  string
	targetBranch        string
}

var killCommand = &cobra.Command{
	Use:   "kill [<branch>]",
	Short: "Removes an obsolete feature branch",
	Long: `Removes an obsolete feature branch

Deletes the current or provided branch from the local and origin repositories.
Does not delete perennial branches nor the main branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := createKillConfig(args, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		stepList, err := createKillStepList(config, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		runState := runstate.New("kill", stepList)
		err = runstate.Execute(runState, prodRepo, nil)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateIsRepository(prodRepo); err != nil {
			return err
		}
		return validateIsConfigured(prodRepo)
	},
}

func createKillConfig(args []string, repo *git.ProdRepo) (killConfig, error) {
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return killConfig{}, err
	}
	result := killConfig{initialBranch: initialBranch}
	if len(args) == 0 {
		result.targetBranch = result.initialBranch
	} else {
		result.targetBranch = args[0]
	}
	if !repo.Config.IsFeatureBranch(result.targetBranch) {
		return result, fmt.Errorf("you can only kill feature branches")
	}
	result.isTargetBranchLocal, err = repo.Silent.HasLocalBranch(result.targetBranch)
	if err != nil {
		return result, err
	}
	if result.isTargetBranchLocal {
		err = dialog.EnsureKnowsParentBranches([]string{result.targetBranch}, repo)
		if err != nil {
			return result, err
		}
		repo.Config.Reload()
	}
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return result, err
	}
	result.isOffline = repo.Config.IsOffline()
	if hasOrigin && !result.isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return result, err
		}
	}
	if result.initialBranch != result.targetBranch {
		hasTargetBranch, err := repo.Silent.HasLocalOrOriginBranch(result.targetBranch)
		if err != nil {
			return result, err
		}
		if !hasTargetBranch {
			return result, fmt.Errorf("there is no branch named %q", result.targetBranch)
		}
	}
	result.hasTrackingBranch, err = repo.Silent.HasTrackingBranch(result.targetBranch)
	if err != nil {
		return result, err
	}
	result.targetBranchParent = repo.Config.ParentBranch(result.targetBranch)
	result.previousBranch, err = repo.Silent.PreviouslyCheckedOutBranch()
	if err != nil {
		return result, err
	}
	result.hasOpenChanges, err = repo.Silent.HasOpenChanges()
	if err != nil {
		return result, err
	}
	result.childBranches = repo.Config.ChildBranches(result.targetBranch)
	result.noPushVerify = !repo.Config.PushVerify()
	return result, nil
}

func createKillStepList(config killConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	switch {
	case config.isTargetBranchLocal:
		if config.hasTrackingBranch && !config.isOffline {
			result.Append(&steps.DeleteOriginBranchStep{BranchName: config.targetBranch, IsTracking: true, NoPushVerify: config.noPushVerify})
		}
		if config.initialBranch == config.targetBranch {
			if config.hasOpenChanges {
				result.Append(&steps.CommitOpenChangesStep{})
			}
			result.Append(&steps.CheckoutBranchStep{BranchName: config.targetBranchParent})
		}
		result.Append(&steps.DeleteLocalBranchStep{BranchName: config.targetBranch, Force: true})
		for _, child := range config.childBranches {
			result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: config.targetBranchParent})
		}
		result.Append(&steps.DeleteParentBranchStep{BranchName: config.targetBranch})
	case !repo.Config.IsOffline():
		result.Append(&steps.DeleteOriginBranchStep{BranchName: config.targetBranch, IsTracking: false, NoPushVerify: config.noPushVerify})
	default:
		return runstate.StepList{}, fmt.Errorf("cannot delete remote branch %q in offline mode", config.targetBranch)
	}
	err := result.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.initialBranch != config.targetBranch && config.targetBranch == config.previousBranch,
	}, repo)
	return result, err
}

func init() {
	RootCmd.AddCommand(killCommand)
}
