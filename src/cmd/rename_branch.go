package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/vm/interpreter"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/git-town/git-town/v11/src/vm/runstate"
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
- confirm with the "--force"/"-f" option
- registers the new perennial branch name in the local Git Town configuration`

func renameBranchCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addForceFlag, readForceFlag := flags.Bool("force", "f", "Force rename of perennial branch", flags.FlagTypeNonPersistent)
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:   "rename-branch [<old_branch_name>] <new_branch_name>",
		Args:  cobra.RangeArgs(1, 2),
		Short: renameBranchDesc,
		Long:  cmdhelpers.Long(renameBranchDesc, renameBranchHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeRenameBranch(args, readDryRunFlag(cmd), readForceFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addForceFlag(&cmd)
	return &cmd
}

func executeRenameBranch(args []string, dryRun, force, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineRenameBranchConfig(args, force, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "rename-branch",
		DryRun:              dryRun,
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          renameBranchProgram(config),
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		FullConfig:              config.FullConfig,
		RunState:                &runState,
		Run:                     repo.Runner,
		Connector:               nil,
		Verbose:                 verbose,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
	})
}

type renameBranchConfig struct {
	*configdomain.FullConfig
	dryRun         bool
	initialBranch  gitdomain.LocalBranchName
	newBranch      gitdomain.LocalBranchName
	oldBranch      gitdomain.BranchInfo
	previousBranch gitdomain.LocalBranchName
}

func determineRenameBranchConfig(args []string, forceFlag bool, repo *execute.OpenRepoResult, dryRun, verbose bool) (*renameBranchConfig, gitdomain.BranchesStatus, gitdomain.StashSize, bool, error) {
	branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSnapshot, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	var oldBranchName gitdomain.LocalBranchName
	var newBranchName gitdomain.LocalBranchName
	if len(args) == 1 {
		oldBranchName = branchesSnapshot.Active
		newBranchName = gitdomain.NewLocalBranchName(args[0])
	} else {
		oldBranchName = gitdomain.NewLocalBranchName(args[0])
		newBranchName = gitdomain.NewLocalBranchName(args[1])
	}
	if repo.Runner.IsMainBranch(oldBranchName) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.RenameMainBranch)
	}
	if !forceFlag {
		if repo.Runner.IsPerennialBranch(oldBranchName) {
			return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.RenameToSameName)
	}
	oldBranch := branchesSnapshot.Branches.FindByLocalName(oldBranchName)
	if oldBranch == nil {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchDoesntExist, oldBranchName)
	}
	if oldBranch.SyncStatus != gitdomain.SyncStatusUpToDate && oldBranch.SyncStatus != gitdomain.SyncStatusLocalOnly {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.RenameBranchNotInSync, oldBranchName)
	}
	if branchesSnapshot.Branches.HasLocalBranch(newBranchName) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, newBranchName)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(newBranchName) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, newBranchName)
	}
	return &renameBranchConfig{
		FullConfig:     &repo.Runner.FullConfig,
		dryRun:         dryRun,
		initialBranch:  branchesSnapshot.Active,
		newBranch:      newBranchName,
		oldBranch:      *oldBranch,
		previousBranch: previousBranch,
	}, branchesSnapshot, stashSnapshot, false, err
}

func renameBranchProgram(config *renameBranchConfig) program.Program {
	result := program.Program{}
	result.Add(&opcode.CreateBranch{Branch: config.newBranch, StartingPoint: config.oldBranch.LocalName.Location()})
	if config.initialBranch == config.oldBranch.LocalName {
		result.Add(&opcode.Checkout{Branch: config.newBranch})
	}
	if !config.dryRun {
		if config.IsPerennialBranch(config.initialBranch) {
			result.Add(&opcode.RemoveFromPerennialBranches{Branch: config.oldBranch.LocalName})
			result.Add(&opcode.AddToPerennialBranches{Branch: config.newBranch})
		} else {
			result.Add(&opcode.DeleteParentBranch{Branch: config.oldBranch.LocalName})
			result.Add(&opcode.SetParent{Branch: config.newBranch, Parent: config.Lineage.Parent(config.oldBranch.LocalName)})
		}
	}
	for _, child := range config.Lineage.Children(config.oldBranch.LocalName) {
		result.Add(&opcode.SetParent{Branch: child, Parent: config.newBranch})
	}
	if config.oldBranch.HasTrackingBranch() && config.IsOnline() {
		result.Add(&opcode.CreateTrackingBranch{Branch: config.newBranch})
		result.Add(&opcode.DeleteTrackingBranch{Branch: config.oldBranch.RemoteName})
	}
	result.Add(&opcode.DeleteLocalBranch{Branch: config.oldBranch.LocalName, Force: false})
	cmdhelpers.Wrap(&result, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             false,
		StashOpenChanges:         false,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch, config.newBranch},
	})
	return result
}
