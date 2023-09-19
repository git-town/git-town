package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/runvm"
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
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, exit, err := determineRenameBranchConfig(args, force, &repo)
	if err != nil || exit {
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
	return runvm.Execute(runvm.ExecuteArgs{
		RunState:  &runState,
		Run:       &repo.Runner,
		Connector: nil,
		Lineage:   config.lineage,
		RootDir:   repo.RootDir,
	})
}

type renameBranchConfig struct {
	branches       domain.Branches
	isOffline      bool
	lineage        config.Lineage
	mainBranch     domain.LocalBranchName
	newBranch      domain.LocalBranchName
	noPushHook     bool
	oldBranch      domain.BranchInfo
	previousBranch domain.LocalBranchName
}

func determineRenameBranchConfig(args []string, forceFlag bool, repo *execute.OpenRepoResult) (*renameBranchConfig, bool, error) {
	lineage := repo.Runner.Config.Lineage()
	branches, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, false, err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	var oldBranchName domain.LocalBranchName
	var newBranchName domain.LocalBranchName
	if len(args) == 1 {
		oldBranchName = branches.Initial
		newBranchName = domain.NewLocalBranchName(args[0])
	} else {
		oldBranchName = domain.NewLocalBranchName(args[0])
		newBranchName = domain.NewLocalBranchName(args[1])
	}
	if repo.Runner.Config.IsMainBranch(oldBranchName) {
		return nil, false, fmt.Errorf(messages.RenameMainBranch)
	}
	if !forceFlag {
		if branches.Types.IsPerennialBranch(oldBranchName) {
			return nil, false, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return nil, false, fmt.Errorf(messages.RenameToSameName)
	}
	oldBranch := branches.All.FindLocalBranch(oldBranchName)
	if oldBranch == nil {
		return nil, false, fmt.Errorf(messages.BranchDoesntExist, oldBranchName)
	}
	if oldBranch.SyncStatus != domain.SyncStatusUpToDate {
		return nil, false, fmt.Errorf(messages.RenameBranchNotInSync, oldBranchName)
	}
	if branches.All.HasLocalBranch(newBranchName) {
		return nil, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, newBranchName)
	}
	if branches.All.HasMatchingRemoteBranchFor(newBranchName) {
		return nil, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, newBranchName)
	}
	return &renameBranchConfig{
		branches:       branches,
		isOffline:      repo.IsOffline,
		lineage:        lineage,
		mainBranch:     mainBranch,
		newBranch:      newBranchName,
		noPushHook:     !pushHook,
		oldBranch:      *oldBranch,
		previousBranch: previousBranch,
	}, false, err
}

func renameBranchStepList(config *renameBranchConfig) (runstate.StepList, error) {
	result := runstate.StepList{}
	result.Append(&steps.CreateBranchStep{Branch: config.newBranch, StartingPoint: config.oldBranch.LocalName.Location()})
	if config.branches.Initial == config.oldBranch.LocalName {
		result.Append(&steps.CheckoutStep{Branch: config.newBranch})
	}
	if config.branches.Types.IsPerennialBranch(config.branches.Initial) {
		result.Append(&steps.RemoveFromPerennialBranchesStep{Branch: config.oldBranch.LocalName})
		result.Append(&steps.AddToPerennialBranchesStep{Branch: config.newBranch})
	} else {
		lineage := config.lineage
		result.Append(&steps.DeleteParentBranchStep{Branch: config.oldBranch.LocalName, Parent: lineage.Parent(config.oldBranch.LocalName)})
		result.Append(&steps.SetParentStep{Branch: config.newBranch, ParentBranch: lineage.Parent(config.oldBranch.LocalName)})
	}
	for _, child := range config.lineage.Children(config.oldBranch.LocalName) {
		result.Append(&steps.SetParentStep{Branch: child, ParentBranch: config.newBranch})
	}
	if config.oldBranch.HasTrackingBranch() && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{Branch: config.newBranch, NoPushHook: config.noPushHook})
		result.Append(&steps.DeleteTrackingBranchStep{Branch: config.oldBranch.LocalName, NoPushHook: false})
	}
	result.Append(&steps.DeleteLocalBranchStep{Branch: config.oldBranch.LocalName, Parent: config.mainBranch.Location(), Force: false})
	err := result.Wrap(runstate.WrapOptions{
		RunInGitRoot:     false,
		StashOpenChanges: false,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return result, err
}
