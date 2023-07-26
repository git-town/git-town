package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
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
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 true,
		HandleUnfinishedState: true,
		OmitBranchNames:       false,
		ValidateIsOnline:      false,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	allBranches, initialBranch, err := execute.LoadBranches(&repo.ProdRunner, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	config, err := determineRenameBranchConfig(args, force, &repo.ProdRunner, allBranches, initialBranch, repo.IsOffline)
	if err != nil {
		return err
	}
	stepList, err := renameBranchStepList(config, &repo.ProdRunner)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "rename-branch",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &repo.ProdRunner, nil, repo.RootDir)
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

func determineRenameBranchConfig(args []string, forceFlag bool, run *git.ProdRunner, allBranches git.BranchesSyncStatus, initialBranch string, isOffline bool) (*renameBranchConfig, error) {
	pushHook, err := run.Config.PushHook()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	var oldBranchName string
	var newBranchName string
	if len(args) == 1 {
		oldBranchName = initialBranch
		newBranchName = args[0]
	} else {
		oldBranchName = args[0]
		newBranchName = args[1]
	}
	if run.Config.IsMainBranch(oldBranchName) {
		return nil, fmt.Errorf(messages.RenameMainBranch)
	}
	if !forceFlag {
		if run.Config.IsPerennialBranch(oldBranchName) {
			return nil, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return nil, fmt.Errorf(messages.RenameToSameName)
	}
	oldBranch := allBranches.Lookup(oldBranchName)
	if oldBranch == nil {
		// TODO: extract these error messages to constants because this one here is reused in several places
		return nil, fmt.Errorf(messages.BranchDoesntExist, oldBranchName)
	}
	if oldBranch.SyncStatus != git.SyncStatusUpToDate {
		return nil, fmt.Errorf(messages.RenameBranchNotInSync, oldBranchName)
	}
	if allBranches.Contains(newBranchName) {
		return nil, fmt.Errorf(messages.BranchAlreadyExists, newBranchName)
	}
	return &renameBranchConfig{
		initialBranch:              initialBranch,
		isInitialBranchPerennial:   run.Config.IsPerennialBranch(initialBranch),
		isOffline:                  isOffline,
		mainBranch:                 mainBranch,
		newBranch:                  newBranchName,
		noPushHook:                 !pushHook,
		oldBranch:                  oldBranchName,
		oldBranchChildren:          run.Config.Lineage().Children(oldBranchName),
		oldBranchHasTrackingBranch: oldBranch.HasTrackingBranch(),
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
