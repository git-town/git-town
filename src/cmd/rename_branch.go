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
	"github.com/git-town/git-town/v14/src/validate"
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
	data, initialBranchesSnapshot, initialStashSize, exit, err := determineRenameBranchData(args, force, repo, dryRun, verbose)
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
		RunProgram:            renameBranchProgram(data),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               nil,
		DialogTestInputs:        &data.dialogTestInputs,
		FinalMessages:           &repo.FinalMessages,
		Frontend:                repo.Frontend,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type renameBranchData struct {
	config           configdomain.ValidatedConfig
	dialogTestInputs components.TestInputs
	dryRun           bool
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	newBranch        gitdomain.LocalBranchName
	oldBranch        gitdomain.BranchInfo
	previousBranch   gitdomain.LocalBranchName
}

func determineRenameBranchData(args []string, forceFlag bool, repo *execute.OpenRepoResult, dryRun, verbose bool) (*renameBranchData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               &repo.Backend,
		Config:                repo.UnvalidatedConfig.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		Frontend:              &repo.Frontend,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	var oldBranchName gitdomain.LocalBranchName
	var newBranchName gitdomain.LocalBranchName
	if len(args) == 1 {
		oldBranchName = branchesSnapshot.Active
		newBranchName = gitdomain.NewLocalBranchName(args[0])
	} else {
		oldBranchName = gitdomain.NewLocalBranchName(args[0])
		newBranchName = gitdomain.NewLocalBranchName(args[1])
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, abort, err := validate.Config(validate.ConfigArgs{
		Backend:            &repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{oldBranchName},
		CommandsCounter:    repo.CommandsCounter,
		ConfigSnapshot:     repo.ConfigSnapshot,
		DialogTestInputs:   dialogTestInputs,
		FinalMessages:      &repo.FinalMessages,
		Frontend:           repo.Frontend,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		RootDir:            repo.RootDir,
		StashSize:          stashSize,
		TestInputs:         &dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
		Verbose:            verbose,
	})
	if err != nil || abort {
		return nil, branchesSnapshot, stashSize, abort, err
	}
	if validatedConfig.Config.IsMainBranch(oldBranchName) {
		return nil, branchesSnapshot, stashSize, false, errors.New(messages.RenameMainBranch)
	}
	if !forceFlag {
		if validatedConfig.Config.IsPerennialBranch(oldBranchName) {
			return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return nil, branchesSnapshot, stashSize, false, errors.New(messages.RenameToSameName)
	}
	oldBranch, hasOldBranch := branchesSnapshot.Branches.FindByLocalName(oldBranchName).Get()
	if !hasOldBranch {
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
	return &renameBranchData{
		config:           validatedConfig.Config,
		dialogTestInputs: dialogTestInputs,
		dryRun:           dryRun,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    branchesSnapshot.Active,
		newBranch:        newBranchName,
		oldBranch:        oldBranch,
		previousBranch:   previousBranch,
	}, branchesSnapshot, stashSize, false, err
}

func renameBranchProgram(data *renameBranchData) program.Program {
	result := program.Program{}
	result.Add(&opcodes.CreateBranch{Branch: data.newBranch, StartingPoint: data.oldBranch.LocalName.Location()})
	if data.initialBranch == data.oldBranch.LocalName {
		result.Add(&opcodes.Checkout{Branch: data.newBranch})
	}
	if !data.dryRun {
		if data.config.IsPerennialBranch(data.initialBranch) {
			result.Add(&opcodes.RemoveFromPerennialBranches{Branch: data.oldBranch.LocalName})
			result.Add(&opcodes.AddToPerennialBranches{Branch: data.newBranch})
		} else {
			result.Add(&opcodes.DeleteParentBranch{Branch: data.oldBranch.LocalName})
			parentBranch, hasParent := data.config.Lineage.Parent(data.oldBranch.LocalName).Get()
			if hasParent {
				result.Add(&opcodes.SetParent{Branch: data.newBranch, Parent: parentBranch})
			}
		}
	}
	for _, child := range data.config.Lineage.Children(data.oldBranch.LocalName) {
		result.Add(&opcodes.SetParent{Branch: child, Parent: data.newBranch})
	}
	if data.oldBranch.HasTrackingBranch() && data.config.IsOnline() {
		result.Add(&opcodes.CreateTrackingBranch{Branch: data.newBranch})
		result.Add(&opcodes.DeleteTrackingBranch{Branch: data.oldBranch.RemoteName})
	}
	result.Add(&opcodes.DeleteLocalBranch{Branch: data.oldBranch.LocalName})
	cmdhelpers.Wrap(&result, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             false,
		StashOpenChanges:         false,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{data.previousBranch, data.newBranch},
	})
	return result
}
