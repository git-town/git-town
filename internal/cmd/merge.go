package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/cli/flags"
	"github.com/git-town/git-town/v18/internal/cli/print"
	"github.com/git-town/git-town/v18/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v18/internal/config"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/execute"
	"github.com/git-town/git-town/v18/internal/forge"
	"github.com/git-town/git-town/v18/internal/forge/forgedomain"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/messages"
	"github.com/git-town/git-town/v18/internal/undo/undoconfig"
	"github.com/git-town/git-town/v18/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v18/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v18/internal/vm/opcodes"
	"github.com/git-town/git-town/v18/internal/vm/optimizer"
	"github.com/git-town/git-town/v18/internal/vm/program"
	"github.com/git-town/git-town/v18/internal/vm/runstate"
	. "github.com/git-town/git-town/v18/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	mergeCmd  = "merge"
	mergeDesc = "Merges the current branch with its parent"
	mergeHelp = `
Merges the current branch with its parent branch.
Both branches must be feature branches.

Consider this branch stack:

main
 \
  branch-1
   \
    branch-2
     \
*     branch-3
       \
        branch-4

We are on the "branch-3" branch. After running "git town merge",
the new "branch-3" branch contains the changes
from the old "branch-2" and "branch-3" branches.

main
 \
  branch-1
   \
*   branch-3
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
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               mergeCmd,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
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
	branchesSnapshot                gitdomain.BranchesSnapshot
	config                          config.ValidatedConfig
	connector                       Option[forgedomain.Connector]
	dialogTestInputs                components.TestInputs
	grandParentBranch               gitdomain.LocalBranchName
	hasOpenChanges                  bool
	initialBranch                   gitdomain.LocalBranchName
	initialBranchFirstCommitMessage Option[gitdomain.CommitMessage]
	initialBranchInfo               gitdomain.BranchInfo
	initialBranchProposal           Option[forgedomain.Proposal]
	initialBranchType               configdomain.BranchType
	offline                         configdomain.Offline
	parentBranch                    gitdomain.LocalBranchName
	parentBranchFirstCommitMessage  Option[gitdomain.CommitMessage]
	parentBranchInfo                gitdomain.BranchInfo
	parentBranchProposal            Option[forgedomain.Proposal]
	parentBranchType                configdomain.BranchType
	prefetchBranchesSnapshot        gitdomain.BranchesSnapshot
	previousBranch                  Option[gitdomain.LocalBranchName]
	remotes                         gitdomain.Remotes
	stashSize                       gitdomain.StashSize
}

func determineMergeData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (mergeData, bool, error) {
	preFetchBranchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return mergeData{}, false, err
	}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return mergeData{}, false, err
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
	grandParentBranch, hasGrandParentBranch := validatedConfig.NormalConfig.Lineage.Parent(parentBranch).Get()
	if !hasGrandParentBranch {
		return mergeData{}, false, fmt.Errorf(messages.MergeNoGrandParent, initialBranch, parentBranch)
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return mergeData{}, false, err
	}
	initialBranchInfo, hasInitialBranchInfo := branchesSnapshot.Branches.FindByLocalName(initialBranch).Get()
	if !hasInitialBranchInfo {
		return mergeData{}, false, fmt.Errorf(messages.BranchInfoNotFound, initialBranch)
	}
	parentBranchInfo, hasParentBranchInfo := branchesSnapshot.Branches.FindByLocalName(parentBranch).Get()
	if !hasParentBranchInfo {
		return mergeData{}, false, fmt.Errorf(messages.BranchInfoNotFound, parentBranch)
	}
	initialBranchFirstCommitMessage, err := repo.Git.FirstCommitMessageInBranch(repo.Backend, initialBranch.BranchName(), parentBranch.BranchName())
	if err != nil {
		return mergeData{}, false, err
	}
	initialBranchType := validatedConfig.BranchType(initialBranch)
	parentBranchType := validatedConfig.BranchType(parentBranch)
	parentBranchFirstCommitMessage, err := repo.Git.FirstCommitMessageInBranch(repo.Backend, parentBranch.BranchName(), grandParentBranch.BranchName())
	if err != nil {
		return mergeData{}, false, err
	}
	initialBranchProposal := None[forgedomain.Proposal]()
	parentBranchProposal := None[forgedomain.Proposal]()
	if connector, hasConnector := connectorOpt.Get(); hasConnector {
		if findProposal, canFindProposal := connector.FindProposalFn().Get(); canFindProposal {
			initialBranchProposal, err = findProposal(initialBranch, parentBranch)
			if err != nil {
				print.Error(err)
			}
			parentBranchProposal, err = findProposal(initialBranch, parentBranch)
			if err != nil {
				print.Error(err)
			}
		}
	}
	return mergeData{
		branchesSnapshot:                branchesSnapshot,
		config:                          validatedConfig,
		connector:                       connectorOpt,
		dialogTestInputs:                dialogTestInputs,
		grandParentBranch:               grandParentBranch,
		hasOpenChanges:                  repoStatus.OpenChanges,
		initialBranch:                   initialBranch,
		initialBranchFirstCommitMessage: initialBranchFirstCommitMessage,
		initialBranchInfo:               *initialBranchInfo,
		initialBranchProposal:           initialBranchProposal,
		initialBranchType:               initialBranchType,
		offline:                         repo.IsOffline,
		parentBranch:                    parentBranch,
		parentBranchFirstCommitMessage:  parentBranchFirstCommitMessage,
		parentBranchInfo:                *parentBranchInfo,
		parentBranchProposal:            parentBranchProposal,
		parentBranchType:                parentBranchType,
		prefetchBranchesSnapshot:        preFetchBranchesSnapshot,
		previousBranch:                  previousBranch,
		remotes:                         remotes,
		stashSize:                       stashSize,
	}, false, err
}

func mergeProgram(data mergeData, dryRun configdomain.DryRun) program.Program {
	prog := NewMutable(&program.Program{})
	if connector, hasConnector := data.connector.Get(); hasConnector && data.offline.IsFalse() {
		initialBranchProposal, hasInitialBranchProposal := data.initialBranchProposal.Get()
		parentBranchProposal, hasParentBranchProposal := data.parentBranchProposal.Get()
		_, connectorCanUpdateSourceBranch := connector.UpdateProposalSourceFn().Get()
		_, connectorCanUpdateTargetBranch := connector.UpdateProposalTargetFn().Get()
		if hasInitialBranchProposal && connectorCanUpdateTargetBranch {
			prog.Value.Add(&opcodes.ProposalUpdateTarget{
				NewBranch:      data.grandParentBranch,
				OldBranch:      data.parentBranch,
				ProposalNumber: initialBranchProposal.Number,
			})
		} else if hasParentBranchProposal && connectorCanUpdateSourceBranch {
			prog.Value.Add(&opcodes.ProposalUpdateSource{
				NewBranch:      data.initialBranch,
				OldBranch:      data.parentBranch,
				ProposalNumber: parentBranchProposal.Number,
			})
		}
	}
	prog.Value.Add(&opcodes.LineageParentSet{
		Branch: data.initialBranch,
		Parent: data.grandParentBranch,
	})
	prog.Value.Add(&opcodes.LineageParentRemove{
		Branch: data.parentBranch,
	})
	prog.Value.Add(&opcodes.BranchLocalDelete{
		Branch: data.parentBranch,
	})
	if data.parentBranchInfo.HasTrackingBranch() && data.offline.IsFalse() {
		prog.Value.Add(&opcodes.BranchTrackingDelete{
			Branch: data.parentBranch.AtRemote(data.config.NormalConfig.DevRemote),
		})
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   dryRun,
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
	inSyncWithParent, err := repo.Git.BranchInSyncWithParent(repo.Backend, data.initialBranch, data.parentBranch)
	if err != nil {
		return err
	}
	if !inSyncWithParent {
		return fmt.Errorf(messages.MergeNotInSyncWithParent, data.initialBranch)
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
		return fmt.Errorf(messages.MergeNotInSyncWithParent, data.parentBranch)
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
