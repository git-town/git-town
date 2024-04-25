package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/spf13/cobra"
)

const renameBranchDesc = "Rename a branch both locally and remotely"

const renameBranchHelp = `
Renames the given branch in the local and origin repository. Aborts if the new branch name already exists or the tracking branch is out of sync.

- creates a branch with the new name
- deletes the old branch

When there is an origin repository:
- syncs the repository

When there is a tracking branch:
- pushes the new branch to the origin repository
- deletes the old branch from the origin repository

When run on a perennial branch:
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
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineRenameBranchConfig(args, force, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               "rename-branch",
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            renameBranchProgram(config),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               nil,
		DialogTestInputs:        &config.dialogTestInputs,
		FullConfig:              config.FullConfig,
		HasOpenChanges:          config.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     repo.Runner,
		RunState:                &runState,
		Verbose:                 verbose,
	})
}

type renameBranchConfig struct {
	*configdomain.FullConfig
	dialogTestInputs components.TestInputs
	dryRun           bool
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	newBranch        gitdomain.LocalBranchName
	oldBranch        gitdomain.BranchInfo
	previousBranch   gitdomain.LocalBranchName
}

func determineRenameBranchConfig(args []string, forceFlag bool, repo *execute.OpenRepoResult, dryRun, verbose bool) (*renameBranchConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Runner.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
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
	if repo.Runner.Config.FullConfig.IsMainBranch(oldBranchName) {
		return nil, branchesSnapshot, stashSize, false, errors.New(messages.RenameMainBranch)
	}
	if !forceFlag {
		if repo.Runner.Config.FullConfig.IsPerennialBranch(oldBranchName) {
			return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return nil, branchesSnapshot, stashSize, false, errors.New(messages.RenameToSameName)
	}
	oldBranch := branchesSnapshot.Branches.FindByLocalName(oldBranchName)
	if oldBranch == nil {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchDoesntExist, oldBranchName)
	}
	if oldBranch.SyncStatus != gitdomain.SyncStatusUpToDate && oldBranch.SyncStatus != gitdomain.SyncStatusLocalOnly {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.RenameBranchNotInSync, oldBranchName)
	}
	if branchesSnapshot.Branches.HasLocalBranch(newBranchName) {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, newBranchName)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(newBranchName) {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, newBranchName)
	}
	return &renameBranchConfig{
		FullConfig:       &repo.Runner.Config.FullConfig,
		dialogTestInputs: dialogTestInputs,
		dryRun:           dryRun,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    branchesSnapshot.Active,
		newBranch:        newBranchName,
		oldBranch:        *oldBranch,
		previousBranch:   previousBranch,
	}, branchesSnapshot, stashSize, false, err
}

func renameBranchProgram(config *renameBranchConfig) program.Program {
	result := program.Program{}
	result.Add(&opcodes.CreateBranch{Branch: config.newBranch, StartingPoint: config.oldBranch.LocalName.Location()})
	if config.initialBranch == config.oldBranch.LocalName {
		result.Add(&opcodes.Checkout{Branch: config.newBranch})
	}
	if !config.dryRun {
		if config.IsPerennialBranch(config.initialBranch) {
			result.Add(&opcodes.RemoveFromPerennialBranches{Branch: config.oldBranch.LocalName})
			result.Add(&opcodes.AddToPerennialBranches{Branch: config.newBranch})
		} else {
			result.Add(&opcodes.DeleteParentBranch{Branch: config.oldBranch.LocalName})
			result.Add(&opcodes.SetParent{Branch: config.newBranch, Parent: config.Lineage.Parent(config.oldBranch.LocalName)})
		}
	}
	for _, child := range config.Lineage.Children(config.oldBranch.LocalName) {
		result.Add(&opcodes.SetParent{Branch: child, Parent: config.newBranch})
	}
	if config.oldBranch.HasTrackingBranch() && config.IsOnline() {
		result.Add(&opcodes.CreateTrackingBranch{Branch: config.newBranch})
		result.Add(&opcodes.DeleteTrackingBranch{Branch: config.oldBranch.RemoteName})
	}
	result.Add(&opcodes.DeleteLocalBranch{Branch: config.oldBranch.LocalName})
	cmdhelpers.Wrap(&result, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             false,
		StashOpenChanges:         false,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch, config.newBranch},
	})
	return result
}
