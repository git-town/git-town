package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
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
	config, err := determineRenameBranchConfig(args, force, &repo.Runner, repo.IsOffline)
	if err != nil {
		return err
	}
	stepList, err := renameBranchStepList(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "rename-branch",
		RunStepList: stepList,
	}
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &runState,
		Run:       &repo.Runner,
		Connector: nil,
		RootDir:   repo.RootDir,
	})
}

type renameBranchConfig struct {
	branchDurations config.BranchDurations
	initialBranch   domain.LocalBranchName
	isOffline       bool
	lineage         config.Lineage
	mainBranch      domain.LocalBranchName
	newBranch       domain.LocalBranchName
	noPushHook      bool
	oldBranch       git.BranchSyncStatus
	previousBranch  domain.LocalBranchName
}

func determineRenameBranchConfig(args []string, forceFlag bool, run *git.ProdRunner, isOffline bool) (*renameBranchConfig, error) {
	branches, err := execute.LoadBranches(run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return nil, err
	}
	previousBranch := run.Backend.PreviouslyCheckedOutBranch()
	pushHook, err := run.Config.PushHook()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	var oldBranchName domain.LocalBranchName
	var newBranchName domain.LocalBranchName
	if len(args) == 1 {
		oldBranchName = branches.Initial
		newBranchName = domain.NewLocalBranchName(args[0])
	} else {
		oldBranchName = domain.NewLocalBranchName(args[0])
		newBranchName = domain.NewLocalBranchName(args[1])
	}
	if run.Config.IsMainBranch(oldBranchName) {
		return nil, fmt.Errorf(messages.RenameMainBranch)
	}
	if !forceFlag {
		if branches.Durations.IsPerennialBranch(oldBranchName) {
			return nil, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return nil, fmt.Errorf(messages.RenameToSameName)
	}
	oldBranch := branches.All.Lookup(oldBranchName)
	if oldBranch == nil {
		// TODO: extract these error messages to constants because this one here is reused in several places
		return nil, fmt.Errorf(messages.BranchDoesntExist, oldBranchName)
	}
	if oldBranch.SyncStatus != git.SyncStatusUpToDate {
		return nil, fmt.Errorf(messages.RenameBranchNotInSync, oldBranchName)
	}
	if branches.All.ContainsLocalBranch(newBranchName) {
		return nil, fmt.Errorf(messages.BranchAlreadyExistsLocally, newBranchName)
	}
	if branches.All.KnowsRemoteBranch(newBranchName.RemoteName()) {
		return nil, fmt.Errorf(messages.BranchAlreadyExistsRemotely, newBranchName)
	}
	lineage := run.Config.Lineage()
	return &renameBranchConfig{
		branchDurations: branches.Durations,
		initialBranch:   branches.Initial,
		isOffline:       isOffline,
		lineage:         lineage,
		mainBranch:      mainBranch,
		newBranch:       newBranchName,
		noPushHook:      !pushHook,
		oldBranch:       *oldBranch,
		previousBranch:  previousBranch,
	}, err
}

func renameBranchStepList(config *renameBranchConfig) (runstate.StepList, error) {
	result := runstate.StepList{}
	result.Append(&steps.CreateBranchStep{Branch: config.newBranch, StartingPoint: config.oldBranch.Name.Location})
	if config.initialBranch == config.oldBranch.Name {
		result.Append(&steps.CheckoutStep{Branch: config.newBranch})
	}
	if config.branchDurations.IsPerennialBranch(config.initialBranch) {
		result.Append(&steps.RemoveFromPerennialBranchesStep{Branch: config.oldBranch.Name})
		result.Append(&steps.AddToPerennialBranchesStep{Branch: config.newBranch})
	} else {
		lineage := config.lineage
		result.Append(&steps.DeleteParentBranchStep{Branch: config.oldBranch.Name, Parent: lineage.Parent(config.oldBranch.Name)})
		result.Append(&steps.SetParentStep{Branch: config.newBranch, ParentBranch: lineage.Parent(config.oldBranch.Name)})
	}
	for _, child := range config.lineage.Children(config.oldBranch.Name) {
		result.Append(&steps.SetParentStep{Branch: child, ParentBranch: config.newBranch})
	}
	if config.oldBranch.HasTrackingBranch() && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{Branch: config.newBranch, NoPushHook: config.noPushHook})
		result.Append(&steps.DeleteOriginBranchStep{Branch: config.oldBranch.Name, IsTracking: true})
	}
	result.Append(&steps.DeleteLocalBranchStep{Branch: config.oldBranch.Name, Parent: config.mainBranch.Location})
	err := result.Wrap(runstate.WrapOptions{
		RunInGitRoot:     false,
		StashOpenChanges: false,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranch,
		PreviousBranch:   config.previousBranch,
	})
	return result, err
}
