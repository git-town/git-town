package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/cmd/sync"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	"github.com/git-town/git-town/v16/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const mergeCmd = "merge"

const mergeDesc = "Merges the current branch with its parent"

const mergeHelp = `
Merges the current branch with its parent branch.
Both branches must be feature branches.
`

func mergeCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     mergeCmd,
		Args:    cobra.NoArgs,
		GroupID: "stack", // TODO: extract into a constant
		Short:   mergeDesc,
		Long:    cmdhelpers.Long(mergeDesc, mergeHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeMerge(readDryRunFlag(cmd), readVerboseFlag(cmd))
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
	err = validateMergeData(data)
	if err != nil {
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
		HasOpenChanges:          false,
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
	branchesSnapshot         gitdomain.BranchesSnapshot
	branchesToSync           configdomain.BranchesToSync
	config                   config.ValidatedConfig
	connector                Option[hostingdomain.Connector]
	dialogTestInputs         components.TestInputs
	grandParentBranch        gitdomain.LocalBranchName
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	initialBranchInfo        gitdomain.BranchInfo
	offline                  configdomain.Offline
	parentBranch             gitdomain.LocalBranchName
	parentBranchInfo         gitdomain.BranchInfo
	prefetchBranchesSnapshot gitdomain.BranchesSnapshot
	previousBranch           Option[gitdomain.LocalBranchName]
	remotes                  gitdomain.Remotes
	stashSize                gitdomain.StashSize
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
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, gitdomain.RemoteOrigin, print.Logger{})
	if err != nil {
		return mergeData{}, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		Connector:          connector,
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
	branchNamesToSync := gitdomain.LocalBranchNames{parentBranch, initialBranch}
	branchesToSync, err := sync.BranchesToSync(branchNamesToSync, branchesSnapshot, repo, validatedConfig.ValidatedConfigData.MainBranch)
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
	return mergeData{
		branchesSnapshot:         branchesSnapshot,
		branchesToSync:           branchesToSync,
		config:                   validatedConfig,
		connector:                connector,
		dialogTestInputs:         dialogTestInputs,
		grandParentBranch:        grandParentBranch,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		initialBranchInfo:        *initialBranchInfo,
		offline:                  repo.IsOffline,
		parentBranch:             parentBranch,
		parentBranchInfo:         *parentBranchInfo,
		prefetchBranchesSnapshot: preFetchBranchesSnapshot,
		previousBranch:           previousBranch,
		remotes:                  remotes,
		stashSize:                stashSize,
	}, false, err
}

func mergeProgram(data mergeData, dryRun configdomain.DryRun) program.Program {
	prog := NewMutable(&program.Program{})
	if data.remotes.HasOrigin() && data.parentBranchInfo.HasTrackingBranch() {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.parentBranch})
		parentBranchSyncInfo := data.branchesToSync.FindByBranch(data.parentBranch).GetOrPanic()
		sync.FeatureTrackingBranchProgram(data.parentBranch.AtRemote(
			gitdomain.RemoteOrigin),
			data.config.NormalConfig.SyncFeatureStrategy.SyncStrategy(),
			sync.FeatureTrackingArgs{
				FirstCommitMessage: parentBranchSyncInfo.FirstCommitMessage,
				LocalName:          data.parentBranch,
				Offline:            data.offline,
				Program:            prog,
				PushBranches:       true,
			})
	}
	initialBranchSyncInfo := data.branchesToSync.FindByBranch(data.initialBranch).GetOrPanic()
	sync.BranchProgram(data.initialBranch, data.initialBranchInfo, initialBranchSyncInfo.FirstCommitMessage, sync.BranchProgramArgs{
		BranchInfos:         data.branchesSnapshot.Branches,
		Config:              data.config,
		InitialBranch:       data.initialBranch,
		PrefetchBranchInfos: data.prefetchBranchesSnapshot.Branches,
		Program:             prog,
		PushBranches:        configdomain.PushBranches(data.initialBranchInfo.HasTrackingBranch()),
		Remotes:             data.remotes,
	})
	// update proposals
	// remove the branches
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
			Branch: data.parentBranch.AtRemote(gitdomain.RemoteOrigin),
		})
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return prog.Get()
}

func validateMergeData(data mergeData) error {
	// ensure parent isn't deleted at remote
	parentInfo, hasParent := data.branchesSnapshot.Branches.FindLocalOrRemote(data.parentBranch).Get()
	if !hasParent {
		return fmt.Errorf(messages.BranchInfoNotFound, data.parentBranch)
	}
	if parentInfo.SyncStatus == gitdomain.SyncStatusDeletedAtRemote {
		return fmt.Errorf(messages.BranchDeletedAtRemote, data.parentBranch)
	}
	if parentInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return fmt.Errorf(messages.BranchOtherWorktree, data.parentBranch)
	}
	// ensure branch isn't deleted at remote
	branchInfo, hasBranchInfo := data.branchesSnapshot.Branches.FindLocalOrRemote(data.initialBranch).Get()
	if !hasBranchInfo {
		return fmt.Errorf(messages.BranchInfoNotFound, data.initialBranch)
	}
	if branchInfo.SyncStatus == gitdomain.SyncStatusDeletedAtRemote {
		return fmt.Errorf(messages.BranchDeletedAtRemote, data.initialBranch)
	}
	// ensure parent branch has only one child
	return nil
}
