package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/vm/interpreter"
	"github.com/git-town/git-town/v10/src/vm/opcode"
	"github.com/git-town/git-town/v10/src/vm/program"
	"github.com/git-town/git-town/v10/src/vm/runstate"
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
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addForceFlag, readForceFlag := flags.Bool("force", "f", "Force rename of perennial branch", flags.FlagTypeNonPersistent)
	cmd := cobra.Command{
		Use:   "rename-branch [<old_branch_name>] <new_branch_name>",
		Args:  cobra.RangeArgs(1, 2),
		Short: renameBranchDesc,
		Long:  long(renameBranchDesc, renameBranchHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeRenameBranch(args, readForceFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	addForceFlag(&cmd)
	return &cmd
}

func executeRenameBranch(args []string, force, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineRenameBranchConfig(args, force, repo, verbose)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "rename-branch",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          renameBranchProgram(config),
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               nil,
		Verbose:                 verbose,
		Lineage:                 config.lineage,
		NoPushHook:              config.noPushHook,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
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

func determineRenameBranchConfig(args []string, forceFlag bool, repo *execute.OpenRepoResult, verbose bool) (*renameBranchConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSnapshot, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
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
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.RenameMainBranch)
	}
	if !forceFlag {
		if branches.Types.IsPerennialBranch(oldBranchName) {
			return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.RenameToSameName)
	}
	oldBranch := branches.All.FindByLocalName(oldBranchName)
	if oldBranch == nil {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchDoesntExist, oldBranchName)
	}
	if oldBranch.SyncStatus != domain.SyncStatusUpToDate {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.RenameBranchNotInSync, oldBranchName)
	}
	if branches.All.HasLocalBranch(newBranchName) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, newBranchName)
	}
	if branches.All.HasMatchingTrackingBranchFor(newBranchName) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, newBranchName)
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
	}, branchesSnapshot, stashSnapshot, false, err
}

func renameBranchProgram(config *renameBranchConfig) program.Program {
	result := program.Program{}
	result.Add(&opcode.CreateBranch{Branch: config.newBranch, StartingPoint: config.oldBranch.LocalName.Location()})
	if config.branches.Initial == config.oldBranch.LocalName {
		result.Add(&opcode.Checkout{Branch: config.newBranch})
	}
	if config.branches.Types.IsPerennialBranch(config.branches.Initial) {
		result.Add(&opcode.RemoveFromPerennialBranches{Branch: config.oldBranch.LocalName})
		result.Add(&opcode.AddToPerennialBranches{Branch: config.newBranch})
	} else {
		lineage := config.lineage
		result.Add(&opcode.DeleteParentBranch{Branch: config.oldBranch.LocalName})
		result.Add(&opcode.SetParent{Branch: config.newBranch, Parent: lineage.Parent(config.oldBranch.LocalName)})
	}
	for _, child := range config.lineage.Children(config.oldBranch.LocalName) {
		result.Add(&opcode.SetParent{Branch: child, Parent: config.newBranch})
	}
	if config.oldBranch.HasTrackingBranch() && !config.isOffline {
		result.Add(&opcode.CreateTrackingBranch{Branch: config.newBranch, NoPushHook: config.noPushHook})
		result.Add(&opcode.DeleteTrackingBranch{Branch: config.oldBranch.RemoteName})
	}
	result.Add(&opcode.DeleteLocalBranch{Branch: config.oldBranch.LocalName, Force: false})
	wrap(&result, wrapOptions{
		RunInGitRoot:     false,
		StashOpenChanges: false,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return result
}
