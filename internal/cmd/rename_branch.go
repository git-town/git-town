package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/cli/flags"
	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v15/internal/config"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/execute"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/git-town/git-town/v15/internal/undo/undoconfig"
	"github.com/git-town/git-town/v15/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v15/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
	"github.com/git-town/git-town/v15/internal/vm/runstate"
	. "github.com/git-town/git-town/v15/pkg/prelude"
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
	addForceFlag, readForceFlag := flags.Force("force rename of perennial branch")
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

func executeRenameBranch(args []string, dryRun configdomain.DryRun, force configdomain.Force, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineRenameBranchData(args, force, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := renameBranchProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "rename-branch",
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               None[hostingdomain.Connector](),
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type renameBranchData struct {
	branchesSnapshot gitdomain.BranchesSnapshot
	config           config.ValidatedConfig
	dialogTestInputs components.TestInputs
	dryRun           configdomain.DryRun
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	newBranch        gitdomain.LocalBranchName
	oldBranch        gitdomain.BranchInfo
	previousBranch   Option[gitdomain.LocalBranchName]
	stashSize        gitdomain.StashSize
}

func determineRenameBranchData(args []string, force configdomain.Force, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data renameBranchData, exit bool, err error) {
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	var oldBranchName gitdomain.LocalBranchName
	var newBranchName gitdomain.LocalBranchName
	if len(args) == 1 {
		oldBranchName = initialBranch
		newBranchName = gitdomain.NewLocalBranchName(args[0])
	} else {
		oldBranchName = gitdomain.NewLocalBranchName(args[0])
		newBranchName = gitdomain.NewLocalBranchName(args[1])
	}
	oldBranch, hasOldBranch := branchesSnapshot.Branches.FindByLocalName(oldBranchName).Get()
	if !hasOldBranch {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, oldBranchName)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{oldBranchName},
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return data, exit, err
	}
	if validatedConfig.Config.IsMainBranch(oldBranchName) {
		return data, false, errors.New(messages.RenameMainBranch)
	}
	if force.IsFalse() {
		if validatedConfig.Config.IsPerennialBranch(oldBranchName) {
			return data, false, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return data, false, errors.New(messages.RenameToSameName)
	}
	if oldBranch.SyncStatus != gitdomain.SyncStatusUpToDate && oldBranch.SyncStatus != gitdomain.SyncStatusLocalOnly {
		return data, false, fmt.Errorf(messages.RenameBranchNotInSync, oldBranchName)
	}
	if branchesSnapshot.Branches.HasLocalBranch(newBranchName) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, newBranchName)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(newBranchName) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, newBranchName)
	}
	return renameBranchData{
		branchesSnapshot: branchesSnapshot,
		config:           validatedConfig,
		dialogTestInputs: dialogTestInputs,
		dryRun:           dryRun,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    initialBranch,
		newBranch:        newBranchName,
		oldBranch:        *oldBranch,
		previousBranch:   previousBranch,
		stashSize:        stashSize,
	}, false, err
}

func renameBranchProgram(data renameBranchData) program.Program {
	result := NewMutable(&program.Program{})
	if oldLocalBranch, hasOldLocalBranch := data.oldBranch.LocalName.Get(); hasOldLocalBranch {
		result.Value.Add(&opcodes.CreateBranch{Branch: data.newBranch, StartingPoint: oldLocalBranch.Location()})
		if data.initialBranch == oldLocalBranch {
			result.Value.Add(&opcodes.Checkout{Branch: data.newBranch})
		}
		if !data.dryRun {
			if data.config.Config.IsPerennialBranch(data.initialBranch) {
				result.Value.Add(&opcodes.RemoveFromPerennialBranches{Branch: oldLocalBranch})
				result.Value.Add(&opcodes.AddToPerennialBranches{Branch: data.newBranch})
			} else {
				result.Value.Add(&opcodes.DeleteParentBranch{Branch: oldLocalBranch})
				parentBranch, hasParent := data.config.Config.Lineage.Parent(oldLocalBranch).Get()
				if hasParent {
					result.Value.Add(&opcodes.SetParent{Branch: data.newBranch, Parent: parentBranch})
				}
			}
		}
		for _, child := range data.config.Config.Lineage.Children(oldLocalBranch) {
			result.Value.Add(&opcodes.SetParent{Branch: child, Parent: data.newBranch})
		}
		if oldTrackingBranch, hasOldTrackingBranch := data.oldBranch.RemoteName.Get(); hasOldTrackingBranch {
			if data.oldBranch.HasTrackingBranch() && data.config.Config.IsOnline() {
				result.Value.Add(&opcodes.CreateTrackingBranch{Branch: data.newBranch})
				result.Value.Add(&opcodes.DeleteTrackingBranch{Branch: oldTrackingBranch})
			}
		}
		result.Value.Add(&opcodes.DeleteLocalBranch{Branch: oldLocalBranch})
		previousBranchCandidates := gitdomain.LocalBranchNames{data.newBranch}
		if previousBranch, hasPrepreviousBranch := data.previousBranch.Get(); hasPrepreviousBranch {
			previousBranchCandidates = append(gitdomain.LocalBranchNames{previousBranch}, previousBranchCandidates...)
		}
		cmdhelpers.Wrap(result, cmdhelpers.WrapOptions{
			DryRun:                   data.dryRun,
			RunInGitRoot:             false,
			StashOpenChanges:         false,
			PreviousBranchCandidates: previousBranchCandidates,
		})
	}
	return result.Get()
}
