package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/spf13/cobra"
)

func renameBranchCommand(repo *git.ProdRepo) *cobra.Command {
	forceFlag := false
	renameBranchCmd := &cobra.Command{
		Use:   "rename-branch [<old_branch_name>] <new_branch_name>",
		Short: "Renames a branch both locally and remotely",
		Long: `Renames a branch both locally and remotely

Renames the given branch in the local and origin repository.
Aborts if the new branch name already exists or the tracking branch is out of sync.

- creates a branch with the new name
- deletes the old branch

When there is an origin repository
- syncs the repository

When there is a tracking branch
- pushes the new branch to the origin repository
- deletes the old branch from the origin repository

When run on a perennial branch
- confirm with the "-f" option
- registers the new perennial branch name in the local Git Town configuration`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := determineRenameBranchConfig(args, forceFlag, repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := renameBranchStepList(config, repo)
			if err != nil {
				cli.Exit(err)
			}
			runState := runstate.New("rename-branch", stepList)
			err = runstate.Execute(runState, repo, nil)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.RangeArgs(1, 2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
	}
	renameBranchCmd.Flags().BoolVar(&forceFlag, "force", false, "Force rename of perennial branch")
	return renameBranchCmd
}

type renameBranchConfig struct {
	initialBranch              string
	isInitialBranchPerennial   bool
	isOffline                  bool
	newBranch                  string
	noPushHook                 bool
	oldBranchChildren          []string
	oldBranchHasTrackingBranch bool
	oldBranch                  string
}

func determineRenameBranchConfig(args []string, forceFlag bool, repo *git.ProdRepo) (*renameBranchConfig, error) {
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return nil, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return nil, err
	}
	pushHook, err := repo.Config.PushHook()
	if err != nil {
		return nil, err
	}
	var oldBranch string
	var newBranch string
	if len(args) == 1 {
		oldBranch = initialBranch
		newBranch = args[0]
	} else {
		oldBranch = args[0]
		newBranch = args[1]
	}
	if repo.Config.IsMainBranch(oldBranch) {
		return nil, fmt.Errorf("the main branch cannot be renamed")
	}
	if !forceFlag {
		if repo.Config.IsPerennialBranch(oldBranch) {
			return nil, fmt.Errorf("%q is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'", oldBranch)
		}
	}
	if oldBranch == newBranch {
		cli.Exit(fmt.Errorf("cannot rename branch to current name"))
	}
	if !isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return nil, err
		}
	}
	hasOldBranch, err := repo.Silent.HasLocalBranch(oldBranch)
	if err != nil {
		return nil, err
	}
	if !hasOldBranch {
		return nil, fmt.Errorf("there is no branch named %q", oldBranch)
	}
	isBranchInSync, err := repo.Silent.IsBranchInSync(oldBranch)
	if err != nil {
		return nil, err
	}
	if !isBranchInSync {
		return nil, fmt.Errorf("%q is not in sync with its tracking branch, please sync the branches before renaming", oldBranch)
	}
	hasNewBranch, err := repo.Silent.HasLocalOrOriginBranch(newBranch)
	if err != nil {
		return nil, err
	}
	if hasNewBranch {
		return nil, fmt.Errorf("a branch named %q already exists", newBranch)
	}
	oldBranchHasTrackingBranch, err := repo.Silent.HasTrackingBranch(oldBranch)
	if err != nil {
		return nil, err
	}
	return &renameBranchConfig{
		initialBranch:              initialBranch,
		isInitialBranchPerennial:   repo.Config.IsPerennialBranch(initialBranch),
		isOffline:                  isOffline,
		newBranch:                  newBranch,
		noPushHook:                 !pushHook,
		oldBranch:                  oldBranch,
		oldBranchChildren:          repo.Config.ChildBranches(oldBranch),
		oldBranchHasTrackingBranch: oldBranchHasTrackingBranch,
	}, err
}

func renameBranchStepList(config *renameBranchConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	result.Append(&steps.CreateBranchStep{Branch: config.newBranch, StartingPoint: config.oldBranch})
	if config.initialBranch == config.oldBranch {
		result.Append(&steps.CheckoutStep{Branch: config.newBranch})
	}
	if config.isInitialBranchPerennial {
		result.Append(&steps.RemoveFromPerennialBranchesStep{Branch: config.oldBranch})
		result.Append(&steps.AddToPerennialBranchesStep{Branch: config.newBranch})
	} else {
		result.Append(&steps.DeleteParentBranchStep{Branch: config.oldBranch})
		result.Append(&steps.SetParentStep{Branch: config.newBranch, ParentBranch: repo.Config.ParentBranch(config.oldBranch)})
	}
	for _, child := range config.oldBranchChildren {
		result.Append(&steps.SetParentStep{Branch: child, ParentBranch: config.newBranch})
	}
	if config.oldBranchHasTrackingBranch && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{Branch: config.newBranch, NoPushHook: config.noPushHook})
		result.Append(&steps.DeleteOriginBranchStep{Branch: config.oldBranch, IsTracking: true})
	}
	result.Append(&steps.DeleteLocalBranchStep{Branch: config.oldBranch})
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: false, StashOpenChanges: false}, repo)
	return result, err
}
