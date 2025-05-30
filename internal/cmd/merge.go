package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/validate"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/optimizer"
	"github.com/git-town/git-town/v21/internal/vm/program"
	"github.com/git-town/git-town/v21/internal/vm/vmstate"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	mergeCmd  = "merge"
	mergeDesc = "Merges the current branch with its parent"
	mergeHelp = `
Merges the current branch with its parent branch.
Both branches must be feature branches.

Consider this stack:

main
 \
  branch-1
   \
    branch-2
     \
*     branch-3
       \
        branch-4

We are on the "branch-3" branch.
After running "git town merge",
the new "branch-3" branch contains the changes
from the old "branch-2" and "branch-3" branches.

main
 \
  branch-1
   \
*   branch-2
     \
      branch-4
`
)

func mergeCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     mergeCmd,
		Args:    cobra.NoArgs,
		GroupID: cmdhelpers.GroupIDStack,
		Short:   mergeDesc,
		Long:    cmdhelpers.Long(mergeDesc, mergeHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeMerge(dryRun, verbose)
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeMerge(dryRun configdomain.DryRun, verbose configdomain.Verbose) error {
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
	data, exit, err := determineMergeData(repo, verbose)
	if err != nil || exit {
		return err
	}
	if err = validateMergeData(repo, data); err != nil {
		return err
	}
	runProgram := mergeProgram(data, dryRun)
	runState := vmstate.Data{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               mergeCmd,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               None[forgedomain.Connector](),
		Detached:                true,
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

type mergeData struct {
	branchInfosLastRun Option[gitdomain.BranchInfos]
	branchesSnapshot   gitdomain.BranchesSnapshot
	config             config.ValidatedConfig
	dialogTestInputs   components.TestInputs
	hasOpenChanges     bool
	initialBranch      gitdomain.LocalBranchName
	initialBranchInfo  gitdomain.BranchInfo
	initialBranchSHA   gitdomain.SHA
	initialBranchType  configdomain.BranchType
	offline            configdomain.Offline
	parentBranch       gitdomain.LocalBranchName
	parentBranchInfo   gitdomain.BranchInfo
	parentBranchSHA    gitdomain.SHA
	parentBranchType   configdomain.BranchType
	previousBranch     Option[gitdomain.LocalBranchName]
	stashSize          gitdomain.StashSize
}

func determineMergeData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (mergeData, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return mergeData{}, false, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Detached:              true,
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
		return mergeData{}, exit, err
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return mergeData{}, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	connectorOpt, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return mergeData{}, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		Connector:          connectorOpt,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return mergeData{}, exit, err
	}
	parentBranch, hasParentBranch := validatedConfig.NormalConfig.Lineage.Parent(initialBranch).Get()
	if !hasParentBranch {
		return mergeData{}, false, fmt.Errorf(messages.MergeNoParent, initialBranch)
	}
	grandParentBranch := validatedConfig.NormalConfig.Lineage.Parent(parentBranch)
	if grandParentBranch.IsNone() {
		return mergeData{}, false, fmt.Errorf(messages.MergeNoGrandParent, initialBranch, parentBranch)
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	initialBranchInfo, hasInitialBranchInfo := branchesSnapshot.Branches.FindByLocalName(initialBranch).Get()
	if !hasInitialBranchInfo {
		return mergeData{}, false, fmt.Errorf(messages.BranchInfoNotFound, initialBranch)
	}
	initialBranchSHA, hasInitialBranchSHA := initialBranchInfo.LocalSHA.Get()
	if !hasInitialBranchSHA {
		return mergeData{}, false, fmt.Errorf(messages.MergeBranchNotLocal, initialBranch)
	}
	parentBranchInfo, hasParentBranchInfo := branchesSnapshot.Branches.FindByLocalName(parentBranch).Get()
	if !hasParentBranchInfo {
		return mergeData{}, false, fmt.Errorf(messages.BranchInfoNotFound, parentBranch)
	}
	parentBranchSHA, hasParentBranchSHA := parentBranchInfo.LocalSHA.Get()
	if !hasParentBranchSHA {
		return mergeData{}, false, fmt.Errorf(messages.MergeBranchNotLocal, parentBranch)
	}
	initialBranchType := validatedConfig.BranchType(initialBranch)
	parentBranchType := validatedConfig.BranchType(parentBranch)
	return mergeData{
		branchInfosLastRun: branchInfosLastRun,
		branchesSnapshot:   branchesSnapshot,
		config:             validatedConfig,
		dialogTestInputs:   dialogTestInputs,
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		initialBranchInfo:  *initialBranchInfo,
		initialBranchSHA:   initialBranchSHA,
		initialBranchType:  initialBranchType,
		offline:            repo.IsOffline,
		parentBranch:       parentBranch,
		parentBranchInfo:   *parentBranchInfo,
		parentBranchSHA:    parentBranchSHA,
		parentBranchType:   parentBranchType,
		previousBranch:     previousBranch,
		stashSize:          stashSize,
	}, false, err
}

func mergeProgram(data mergeData, dryRun configdomain.DryRun) program.Program {
	prog := NewMutable(&program.Program{})
	// there is no point in updating proposals:
	// If the parent branch has a proposal, it doesn't need to change.
	// The child branch proposal will get closed because the child branch gets deleted,
	// and that's correct because it was from the child branch into the parent branch,
	// and that doesn't make sense anymore because both branches are one now.
	prog.Value.Add(&opcodes.Checkout{Branch: data.parentBranch})
	if data.initialBranchSHA != data.parentBranchSHA {
		prog.Value.Add(&opcodes.BranchLocalSetToSHA{SHA: data.initialBranchSHA})
	}
	prog.Value.Add(&opcodes.LineageParentRemove{
		Branch: data.initialBranch,
	})
	prog.Value.Add(&opcodes.BranchLocalDelete{
		Branch: data.initialBranch,
	})
	if data.parentBranchInfo.RemoteName.IsSome() && data.offline.IsOnline() {
		prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{
			CurrentBranch:   data.parentBranch,
			ForceIfIncludes: true,
		})
	}
	initialTrackingBranch, initialHasTrackingBranch := data.initialBranchInfo.RemoteName.Get()
	if initialHasTrackingBranch && data.offline.IsOnline() {
		prog.Value.Add(&opcodes.BranchTrackingDelete{
			Branch: initialTrackingBranch,
		})
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   dryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return optimizer.Optimize(prog.Immutable())
}

func validateMergeData(repo execute.OpenRepoResult, data mergeData) error {
	if err := verifyBranchType(data.initialBranchType); err != nil {
		return err
	}
	if err := verifyBranchType(data.parentBranchType); err != nil {
		return err
	}
	// ensure all commits on parent branch are contained in the initial branch
	inSyncWithParent, err := repo.Git.BranchInSyncWithParent(repo.Backend, data.initialBranch, data.parentBranch.BranchName())
	if err != nil {
		return err
	}
	if !inSyncWithParent {
		return fmt.Errorf(messages.BranchNotInSyncWithParent, data.initialBranch)
	}
	switch data.initialBranchInfo.SyncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusLocalOnly:
	case gitdomain.SyncStatusAhead, gitdomain.SyncStatusBehind, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusDeletedAtRemote:
		return fmt.Errorf(messages.MergeNotInSyncWithTracking, data.initialBranch)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.BranchOtherWorktree, data.parentBranch)
	case gitdomain.SyncStatusRemoteOnly:
		// safe to ignore, this cannot happen
	}
	switch data.parentBranchInfo.SyncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusLocalOnly, gitdomain.SyncStatusRemoteOnly:
	case gitdomain.SyncStatusAhead, gitdomain.SyncStatusBehind, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusDeletedAtRemote:
		return fmt.Errorf(messages.BranchNotInSyncWithParent, data.parentBranch)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.BranchOtherWorktree, data.parentBranch)
	}
	children := data.config.NormalConfig.Lineage.Children(data.parentBranch)
	if len(children) > 1 {
		return fmt.Errorf("branch %q has more than one child", data.parentBranch)
	}
	return nil
}

func verifyBranchType(branchType configdomain.BranchType) error {
	switch branchType {
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypePerennialBranch:
		return fmt.Errorf(messages.MergeWrongBranchType, branchType)
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
	}
	return nil
}
