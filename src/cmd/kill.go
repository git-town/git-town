package cmd

import (
	"fmt"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
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
		config, err := getKillConfig(args, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		stepList, err := getKillStepList(config, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		runState := steps.NewRunState("kill", stepList)
		err = steps.Run(runState, prodRepo, nil)
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

// nolint: funlen
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
		return result, fmt.Errorf("you can only kill feature branches")
	}
	result.isTargetBranchLocal, err = repo.Silent.HasLocalBranch(result.targetBranch)
	if err != nil {
		return result, err
	}
	if result.isTargetBranchLocal {
		err = prompt.EnsureKnowsParentBranches([]string{result.targetBranch}, repo)
		if err != nil {
			return result, err
		}
		repo.Configuration.Reload()
	}
	hasOrigin, err := repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	result.isOffline = repo.IsOffline()
	if hasOrigin && !result.isOffline {
		err := repo.Logging.Fetch()
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
			return result, fmt.Errorf("there is no branch named %q", result.targetBranch)
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

func getKillStepList(config killConfig, repo *git.ProdRepo) (result steps.StepList, err error) {
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
		return result, fmt.Errorf("cannot delete remote branch %q in offline mode", config.targetBranch)
	}
	err = result.Wrap(steps.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.initialBranch != config.targetBranch && config.targetBranch == config.previousBranch,
	}, repo)
	return result, err
}

func init() {
	RootCmd.AddCommand(killCommand)
}
