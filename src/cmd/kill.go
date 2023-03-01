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

func killCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "kill [<branch>]",
		Short: "Removes an obsolete feature branch",
		Long: `Removes an obsolete feature branch

Deletes the current or provided branch from the local and origin repositories.
Does not delete perennial branches nor the main branch.`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := determineKillConfig(args, repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := killStepList(config, repo)
			if err != nil {
				cli.Exit(err)
			}
			runState := runstate.New("kill", stepList)
			err = runstate.Execute(runState, repo, nil)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
	}
}

type killConfig struct {
	childBranches       []string
	hasOpenChanges      bool
	hasTrackingBranch   bool
	initialBranch       string
	isOffline           bool
	isTargetBranchLocal bool
	noPushHook          bool
	previousBranch      string
	targetBranchParent  string
	targetBranch        string
}

func determineKillConfig(args []string, repo *git.ProdRepo) (*killConfig, error) {
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return nil, err
	}
	var targetBranch string
	if len(args) > 0 {
		targetBranch = args[0]
	} else {
		targetBranch = initialBranch
	}
	if !repo.Config.IsFeatureBranch(targetBranch) {
		return nil, fmt.Errorf("you can only kill feature branches")
	}
	isTargetBranchLocal, err := repo.Silent.HasLocalBranch(targetBranch)
	if err != nil {
		return nil, err
	}
	if isTargetBranchLocal {
		parentDialog := dialog.ParentBranches{}
		err = parentDialog.EnsureKnowsParentBranches([]string{targetBranch}, repo)
		if err != nil {
			return nil, err
		}
		repo.Config.Reload()
	}
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return nil, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return nil, err
	}
	if hasOrigin && !isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return nil, err
		}
	}
	if initialBranch != targetBranch {
		hasTargetBranch, err := repo.Silent.HasLocalOrOriginBranch(targetBranch)
		if err != nil {
			return nil, err
		}
		if !hasTargetBranch {
			return nil, fmt.Errorf("there is no branch named %q", targetBranch)
		}
	}
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(targetBranch)
	if err != nil {
		return nil, err
	}
	previousBranch, err := repo.Silent.PreviouslyCheckedOutBranch()
	if err != nil {
		return nil, err
	}
	hasOpenChanges, err := repo.Silent.HasOpenChanges()
	if err != nil {
		return nil, err
	}
	pushHook, err := repo.Config.PushHook()
	if err != nil {
		return nil, err
	}
	return &killConfig{
		childBranches:       repo.Config.ChildBranches(targetBranch),
		hasOpenChanges:      hasOpenChanges,
		hasTrackingBranch:   hasTrackingBranch,
		initialBranch:       initialBranch,
		isOffline:           isOffline,
		isTargetBranchLocal: isTargetBranchLocal,
		noPushHook:          !pushHook,
		previousBranch:      previousBranch,
		targetBranch:        targetBranch,
		targetBranchParent:  repo.Config.ParentBranch(targetBranch),
	}, nil
}

func killStepList(config *killConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	switch {
	case config.isTargetBranchLocal:
		if config.hasTrackingBranch && !config.isOffline {
			result.Append(&steps.DeleteOriginBranchStep{Branch: config.targetBranch, IsTracking: true, NoPushHook: config.noPushHook})
		}
		if config.initialBranch == config.targetBranch {
			if config.hasOpenChanges {
				result.Append(&steps.CommitOpenChangesStep{})
			}
			result.Append(&steps.CheckoutBranchStep{Branch: config.targetBranchParent})
		}
		result.Append(&steps.DeleteLocalBranchStep{Branch: config.targetBranch, Force: true})
		for _, child := range config.childBranches {
			result.Append(&steps.SetParentStep{Branch: child, ParentBranch: config.targetBranchParent})
		}
		result.Append(&steps.DeleteParentBranchStep{Branch: config.targetBranch})
	case !config.isOffline:
		result.Append(&steps.DeleteOriginBranchStep{Branch: config.targetBranch, IsTracking: false, NoPushHook: config.noPushHook})
	default:
		return runstate.StepList{}, fmt.Errorf("cannot delete remote branch %q in offline mode", config.targetBranch)
	}
	err := result.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.initialBranch != config.targetBranch && config.targetBranch == config.previousBranch,
	}, repo)
	return result, err
}
