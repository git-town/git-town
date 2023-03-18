package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func killCommand() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:   "kill [<branch>]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Removes an obsolete feature branch",
		Long: `Removes an obsolete feature branch

Deletes the current or provided branch from the local and origin repositories.
Does not delete perennial branches nor the main branch.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runKill(debug, args)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runKill(debug bool, args []string) error {
	repo, err := Repo(RepoArgs{
		printBranchNames:     false,
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
		validateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	config, err := determineKillConfig(args, &repo)
	if err != nil {
		return err
	}
	stepList, err := killStepList(config, &repo)
	if err != nil {
		return err
	}
	runState := runstate.New("kill", stepList)
	return runstate.Execute(runState, &repo, nil)
}

type killConfig struct {
	childBranches       []string
	hasOpenChanges      bool
	hasTrackingBranch   bool
	initialBranch       string
	isOffline           bool
	isTargetBranchLocal bool
	mainBranch          string
	noPushHook          bool
	previousBranch      string
	targetBranchParent  string
	targetBranch        string
}

func determineKillConfig(args []string, repo *git.PublicRepo) (*killConfig, error) {
	initialBranch, err := repo.CurrentBranch()
	if err != nil {
		return nil, err
	}
	mainBranch := repo.Config.MainBranch()
	var targetBranch string
	if len(args) > 0 {
		targetBranch = args[0]
	} else {
		targetBranch = initialBranch
	}
	if !repo.Config.IsFeatureBranch(targetBranch) {
		return nil, fmt.Errorf("you can only kill feature branches")
	}
	isTargetBranchLocal, err := repo.HasLocalBranch(targetBranch)
	if err != nil {
		return nil, err
	}
	if isTargetBranchLocal {
		err = validate.KnowsBranchAncestry(targetBranch, mainBranch, repo)
		if err != nil {
			return nil, err
		}
		repo.Config.Reload()
	}
	hasOrigin, err := repo.HasOrigin()
	if err != nil {
		return nil, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return nil, err
	}
	if hasOrigin && !isOffline {
		err := repo.Fetch()
		if err != nil {
			return nil, err
		}
	}
	if initialBranch != targetBranch {
		hasTargetBranch, err := repo.HasLocalOrOriginBranch(targetBranch, mainBranch)
		if err != nil {
			return nil, err
		}
		if !hasTargetBranch {
			return nil, fmt.Errorf("there is no branch named %q", targetBranch)
		}
	}
	hasTrackingBranch, err := repo.HasTrackingBranch(targetBranch)
	if err != nil {
		return nil, err
	}
	previousBranch, err := repo.PreviouslyCheckedOutBranch()
	if err != nil {
		return nil, err
	}
	hasOpenChanges, err := repo.HasOpenChanges()
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
		mainBranch:          mainBranch,
		noPushHook:          !pushHook,
		previousBranch:      previousBranch,
		targetBranch:        targetBranch,
		targetBranchParent:  repo.Config.ParentBranch(targetBranch),
	}, nil
}

func killStepList(config *killConfig, repo *git.PublicRepo) (runstate.StepList, error) {
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
			result.Append(&steps.CheckoutStep{Branch: config.targetBranchParent})
		}
		result.Append(&steps.DeleteLocalBranchStep{Branch: config.targetBranch, Parent: config.mainBranch, Force: true})
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
	}, repo, config.mainBranch)
	return result, err
}
