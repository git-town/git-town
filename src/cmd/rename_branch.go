package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/spf13/cobra"
)

type renameBranchConfig struct {
	initialBranch              string
	isInitialBranchPerennial   bool
	isOffline                  bool
	newBranchName              string
	noPushHook                 bool
	oldBranchChildren          []string
	oldBranchHasTrackingBranch bool
	oldBranchName              string
}

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

func determineRenameBranchConfig(args []string, forceFlag bool, repo *git.ProdRepo) (renameBranchConfig, error) {
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return renameBranchConfig{}, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return renameBranchConfig{}, err
	}
	result := renameBranchConfig{
		initialBranch:            initialBranch,
		isInitialBranchPerennial: repo.Config.IsPerennialBranch(initialBranch),
		isOffline:                isOffline,
	}
	if len(args) == 1 {
		result.oldBranchName = result.initialBranch
		result.newBranchName = args[0]
	} else {
		result.oldBranchName = args[0]
		result.newBranchName = args[1]
	}
	if repo.Config.IsMainBranch(result.oldBranchName) {
		return renameBranchConfig{}, fmt.Errorf("the main branch cannot be renamed")
	}
	if !forceFlag {
		if repo.Config.IsPerennialBranch(result.oldBranchName) {
			return renameBranchConfig{}, fmt.Errorf("%q is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'", result.oldBranchName)
		}
	}
	if result.oldBranchName == result.newBranchName {
		cli.Exit(fmt.Errorf("cannot rename branch to current name"))
	}
	if !result.isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return renameBranchConfig{}, err
		}
	}
	hasOldBranch, err := repo.Silent.HasLocalBranch(result.oldBranchName)
	if err != nil {
		return renameBranchConfig{}, err
	}
	if !hasOldBranch {
		return renameBranchConfig{}, fmt.Errorf("there is no branch named %q", result.oldBranchName)
	}
	isBranchInSync, err := repo.Silent.IsBranchInSync(result.oldBranchName)
	if err != nil {
		return renameBranchConfig{}, err
	}
	if !isBranchInSync {
		return renameBranchConfig{}, fmt.Errorf("%q is not in sync with its tracking branch, please sync the branches before renaming", result.oldBranchName)
	}
	hasNewBranch, err := repo.Silent.HasLocalOrOriginBranch(result.newBranchName)
	if err != nil {
		return renameBranchConfig{}, err
	}
	if hasNewBranch {
		return renameBranchConfig{}, fmt.Errorf("a branch named %q already exists", result.newBranchName)
	}
	pushHook, err := repo.Config.PushHook()
	if err != nil {
		return renameBranchConfig{}, err
	}
	result.noPushHook = !pushHook
	result.oldBranchChildren = repo.Config.ChildBranches(result.oldBranchName)
	result.oldBranchHasTrackingBranch, err = repo.Silent.HasTrackingBranch(result.oldBranchName)
	return result, err
}

func renameBranchStepList(config renameBranchConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	result.Append(&steps.CreateBranchStep{BranchName: config.newBranchName, StartingPoint: config.oldBranchName})
	if config.initialBranch == config.oldBranchName {
		result.Append(&steps.CheckoutBranchStep{BranchName: config.newBranchName})
	}
	if config.isInitialBranchPerennial {
		result.Append(&steps.RemoveFromPerennialBranchesStep{BranchName: config.oldBranchName})
		result.Append(&steps.AddToPerennialBranchesStep{BranchName: config.newBranchName})
	} else {
		result.Append(&steps.DeleteParentBranchStep{BranchName: config.oldBranchName})
		result.Append(&steps.SetParentBranchStep{BranchName: config.newBranchName, ParentBranchName: repo.Config.ParentBranch(config.oldBranchName)})
	}
	for _, child := range config.oldBranchChildren {
		result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: config.newBranchName})
	}
	if config.oldBranchHasTrackingBranch && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{BranchName: config.newBranchName, NoPushHook: config.noPushHook})
		result.Append(&steps.DeleteOriginBranchStep{BranchName: config.oldBranchName, IsTracking: true})
	}
	result.Append(&steps.DeleteLocalBranchStep{BranchName: config.oldBranchName})
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: false, StashOpenChanges: false}, repo)
	return result, err
}
