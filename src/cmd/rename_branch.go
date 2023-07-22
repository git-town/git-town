package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/spf13/cobra"
)

const renameBranchDesc = "Renames a branch both locally and remotely"

const renameBranchHelp = `
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
- registers the new perennial branch name in the local Git Town configuration`

func renameBranchCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	addForceFlag, readForceFlag := flags.Bool("force", "f", "Force rename of perennial branch")
	cmd := cobra.Command{
		Use:   "rename-branch [<old_branch_name>] <new_branch_name>",
		Args:  cobra.RangeArgs(1, 2),
		Short: renameBranchDesc,
		Long:  long(renameBranchDesc, renameBranchHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return renameBranch(args, readForceFlag(cmd), readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	addForceFlag(&cmd)
	return &cmd
}

func renameBranch(args []string, force, debug bool) error {
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: false,
	})
	if err != nil {
		return err
	}
	branchesSyncStatus, initialBranch, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateIsOnline:      false,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineRenameBranchConfig(args, force, &run, branchesSyncStatus, initialBranch)
	if err != nil {
		return err
	}
	stepList, err := renameBranchStepList(config, &run)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "rename-branch",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &run, nil)
}

type renameBranchConfig struct {
	initialBranch              string
	isInitialBranchPerennial   bool
	isOffline                  bool
	mainBranch                 string
	newBranch                  string
	noPushHook                 bool
	oldBranchChildren          []string
	oldBranchHasTrackingBranch bool
	oldBranch                  string
}

func determineRenameBranchConfig(args []string, forceFlag bool, run *git.ProdRunner, branchesSyncStatus git.BranchesSyncStatus, initialBranch string) (*renameBranchConfig, error) {
	isOffline, err := run.Config.IsOffline()
	if err != nil {
		return nil, err
	}
	pushHook, err := run.Config.PushHook()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	var oldBranch string
	var newBranch string
	if len(args) == 1 {
		oldBranch = initialBranch
		newBranch = args[0]
	} else {
		oldBranch = args[0]
		newBranch = args[1]
	}
	if run.Config.IsMainBranch(oldBranch) {
		return nil, fmt.Errorf("the main branch cannot be renamed")
	}
	if !forceFlag {
		if run.Config.IsPerennialBranch(oldBranch) {
			return nil, fmt.Errorf("%q is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'", oldBranch)
		}
	}
	if oldBranch == newBranch {
		return nil, fmt.Errorf("cannot rename branch to current name")
	}
	hasOldBranch, err := run.Backend.HasLocalBranch(oldBranch)
	if err != nil {
		return nil, err
	}
	if !hasOldBranch {
		return nil, fmt.Errorf("there is no branch named %q", oldBranch)
	}
	isBranchInSync, err := run.Backend.IsBranchInSync(oldBranch)
	if err != nil {
		return nil, err
	}
	if !isBranchInSync {
		return nil, fmt.Errorf("%q is not in sync with its tracking branch, please sync the branches before renaming", oldBranch)
	}
	if branchesSyncStatus.Contains(newBranch) {
		return nil, fmt.Errorf("a branch named %q already exists", newBranch)
	}
	oldBranchHasTrackingBranch, err := run.Backend.HasTrackingBranch(oldBranch)
	if err != nil {
		return nil, err
	}
	return &renameBranchConfig{
		initialBranch:              initialBranch,
		isInitialBranchPerennial:   run.Config.IsPerennialBranch(initialBranch),
		isOffline:                  isOffline,
		mainBranch:                 mainBranch,
		newBranch:                  newBranch,
		noPushHook:                 !pushHook,
		oldBranch:                  oldBranch,
		oldBranchChildren:          run.Config.Lineage().Children(oldBranch),
		oldBranchHasTrackingBranch: oldBranchHasTrackingBranch,
	}, err
}

func renameBranchStepList(config *renameBranchConfig, run *git.ProdRunner) (runstate.StepList, error) {
	result := runstate.StepList{}
	result.Append(&steps.CreateBranchStep{Branch: config.newBranch, StartingPoint: config.oldBranch})
	if config.initialBranch == config.oldBranch {
		result.Append(&steps.CheckoutStep{Branch: config.newBranch})
	}
	if config.isInitialBranchPerennial {
		result.Append(&steps.RemoveFromPerennialBranchesStep{Branch: config.oldBranch})
		result.Append(&steps.AddToPerennialBranchesStep{Branch: config.newBranch})
	} else {
		lineage := run.Config.Lineage()
		result.Append(&steps.DeleteParentBranchStep{Branch: config.oldBranch, Parent: lineage.Parent(config.oldBranch)})
		result.Append(&steps.SetParentStep{Branch: config.newBranch, ParentBranch: lineage.Parent(config.oldBranch)})
	}
	for _, child := range config.oldBranchChildren {
		result.Append(&steps.SetParentStep{Branch: child, ParentBranch: config.newBranch})
	}
	if config.oldBranchHasTrackingBranch && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{Branch: config.newBranch, NoPushHook: config.noPushHook})
		result.Append(&steps.DeleteOriginBranchStep{Branch: config.oldBranch, IsTracking: true})
	}
	result.Append(&steps.DeleteLocalBranchStep{Branch: config.oldBranch, Parent: config.mainBranch})
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: false, StashOpenChanges: false}, &run.Backend, config.mainBranch)
	return result, err
}
