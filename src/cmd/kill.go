package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/script"
	"github.com/git-town/git-town/src/steps"
	"github.com/spf13/cobra"
)

type killConfig struct {
	initialBranch       string
	previousBranch      string
	targetBranchParent  string
	targetBranch        string
	childBranches       []string
	isOffline           bool
	isTargetBranchLocal bool
	hasOpenChanges      bool
	hasTrackingBranch   bool
}

var killCommand = &cobra.Command{
	Use:   "kill [<branch>]",
	Short: "Removes an obsolete feature branch",
	Long: `Removes an obsolete feature branch

Deletes the current or provided branch from the local and remote repositories.
Does not delete perennial branches nor the main branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := git.NewProdRepo()
		config, err := getKillConfig(args, repo)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		stepList := getKillStepList(config, repo)
		runState := steps.NewRunState("kill", stepList)
		err = steps.Run(runState, repo)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := git.ValidateIsRepository(); err != nil {
			return err
		}
		return validateIsConfigured()
	},
}

func getKillConfig(args []string, repo *git.ProdRepo) (result killConfig, err error) {
	result.initialBranch, err = repo.Silent.CurrentBranch()
	if err != nil {
		return result, err
	}
	if len(args) == 0 {
		result.targetBranch = result.initialBranch
	} else {
		result.targetBranch = args[0]
	}
	if !repo.IsFeatureBranch(result.targetBranch) {
		return result, fmt.Errorf("branch %q: %w", result.targetBranch, ErrNoFeatureBranch)
	}
	result.isTargetBranchLocal, err = repo.Silent.HasLocalBranch(result.targetBranch)
	if err != nil {
		return result, err
	}
	if result.isTargetBranchLocal {
		prompt.EnsureKnowsParentBranches([]string{result.targetBranch})
		repo.ConfigurationInterface.Reload()
	}
	hasOrigin, err := repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	result.isOffline = repo.IsOffline()
	if hasOrigin && !result.isOffline {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}
	if result.initialBranch != result.targetBranch {
		hasTargetBranch, err := repo.Silent.HasLocalOrRemoteBranch(result.targetBranch)
		if err != nil {
			return result, err
		}
		if !hasTargetBranch {
			return result, fmt.Errorf("branch %q: %w", result.targetBranch, ErrBranchMissing)
		}
	}
	result.hasTrackingBranch, err = repo.Silent.HasTrackingBranch(result.targetBranch)
	if err != nil {
		return result, err
	}
	result.targetBranchParent = repo.GetParentBranch(result.targetBranch)
	result.previousBranch, err = repo.Silent.PreviouslyCheckedOutBranch()
	if err != nil {
		return result, err
	}
	result.hasOpenChanges, err = repo.Silent.HasOpenChanges()
	if err != nil {
		return result, err
	}
	result.childBranches = repo.GetChildBranches(result.targetBranch)
	return result, nil
}

func getKillStepList(config killConfig, repo *git.ProdRepo) (result steps.StepList) {
	switch {
	case config.isTargetBranchLocal:
		if config.hasTrackingBranch && !config.isOffline {
			result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.targetBranch, IsTracking: true})
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
	case !repo.IsOffline():
		result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.targetBranch, IsTracking: false})
	default:
		fmt.Printf("Cannot delete remote branch %q in offline mode", config.targetBranch)
		os.Exit(1)
	}
	result.Wrap(steps.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.initialBranch != config.targetBranch && config.targetBranch == config.previousBranch,
	})
	return result
}

func init() {
	RootCmd.AddCommand(killCommand)
}
