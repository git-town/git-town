package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/cmd/ship"
	"github.com/git-town/git-town/v21/internal/cmd/sync"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/validate"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/optimizer"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	deleteDesc = "Remove an obsolete feature branch"
	deleteHelp = `
Deletes the current or provided branch
and its tracking branch.
Does not delete perennial branches
nor the main branch.

Consider this stack:

main
 \
  branch-1
   \
*   branch-2
     \
      branch-3

We are on the "branch-2" branch.
After running "git town delete"
we end up with this stack,
on the branch that was active
before we switched to "branch-2":

main
 \
  branch-1
   \
    branch-3
`
)

func deleteCommand() *cobra.Command {
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "delete [<branch>]",
		Args:  cobra.MaximumNArgs(1),
		Short: deleteDesc,
		Long:  cmdhelpers.Long(deleteDesc, deleteHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, err1 := readDryRunFlag(cmd)
			verbose, err2 := readVerboseFlag(cmd)
			if err := cmp.Or(err1, err2); err != nil {
				return err
			}
			return executeDelete(args, dryRun, verbose)
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeDelete(args []string, dryRun configdomain.DryRun, verbose configdomain.Verbose) error {
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
	data, exit, err := determineDeleteData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateDeleteData(data)
	if err != nil {
		return err
	}
	runProgram, finalUndoProgram := deleteProgram(repo, data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "delete",
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		FinalUndoProgram:      finalUndoProgram,
		BranchInfosLastRun:    data.branchInfosLastRun,
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
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
		PendingCommand:          None[string](),
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type deleteData struct {
	branchInfosLastRun       Option[gitdomain.BranchInfos]
	branchToDeleteInfo       gitdomain.BranchInfo
	branchToDeleteType       configdomain.BranchType
	branchWhenDone           gitdomain.LocalBranchName
	branchesSnapshot         gitdomain.BranchesSnapshot
	config                   config.ValidatedConfig
	connector                Option[forgedomain.Connector]
	dialogTestInputs         dialogcomponents.TestInputs
	dryRun                   configdomain.DryRun
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	nonExistingBranches      gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	previousBranch           Option[gitdomain.LocalBranchName]
	proposalsOfChildBranches []forgedomain.Proposal
	stashSize                gitdomain.StashSize
}

func determineDeleteData(args []string, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data deleteData, exit dialogdomain.Exit, err error) {
	dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		CodebergToken:        config.CodebergToken,
		ForgeType:            config.ForgeType,
		Frontend:             repo.Frontend,
		GitHubConnectorType:  config.GitHubConnectorType,
		GitHubToken:          config.GitHubToken,
		GitLabConnectorType:  config.GitLabConnectorType,
		GitLabToken:          config.GitLabToken,
		GiteaToken:           config.GiteaToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
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
		return data, exit, err
	}
	branchNameToDelete := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToDelete, hasBranchToDelete := branchesSnapshot.Branches.FindByLocalName(branchNameToDelete).Get()
	if !hasBranchToDelete {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToDelete)
	}
	if branchToDelete.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return data, exit, fmt.Errorf(messages.BranchOtherWorktree, branchNameToDelete)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{},
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
		return data, exit, err
	}
	branchTypeToDelete := validatedConfig.BranchType(branchNameToDelete)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	previousBranchOpt := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	branchWhenDone := determineBranchWhenDone(branchWhenDoneArgs{
		branchNameToDelete: branchNameToDelete,
		branches:           branchesSnapshot.Branches,
		initialBranch:      initialBranch,
		mainBranch:         validatedConfig.ValidatedConfigData.MainBranch,
		previousBranch:     previousBranchOpt,
	})
	proposalsOfChildBranches := ship.LoadProposalsOfChildBranches(ship.LoadProposalsOfChildBranchesArgs{
		ConnectorOpt:               connector,
		Lineage:                    validatedConfig.NormalConfig.Lineage,
		Offline:                    repo.IsOffline,
		OldBranch:                  branchNameToDelete,
		OldBranchHasTrackingBranch: branchToDelete.HasTrackingBranch(),
	})
	lineageBranches := validatedConfig.NormalConfig.Lineage.BranchNames()
	_, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, lineageBranches...)
	return deleteData{
		branchInfosLastRun:       branchInfosLastRun,
		branchToDeleteInfo:       *branchToDelete,
		branchToDeleteType:       branchTypeToDelete,
		branchWhenDone:           branchWhenDone,
		branchesSnapshot:         branchesSnapshot,
		config:                   validatedConfig,
		connector:                connector,
		dialogTestInputs:         dialogTestInputs,
		dryRun:                   dryRun,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		nonExistingBranches:      nonExistingBranches,
		previousBranch:           previousBranchOpt,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
	}, false, nil
}

func deleteProgram(repo execute.OpenRepoResult, data deleteData, finalMessages stringslice.Collector) (runProgram, finalUndoProgram program.Program) {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages, repo.Backend)
	undoProg := NewMutable(&program.Program{})
	switch data.branchToDeleteType {
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
		deleteFeatureBranch(prog, undoProg, data)
	case
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeContributionBranch:
		deleteLocalBranch(prog, undoProg, data)
	case
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		panic(fmt.Sprintf("this branch type should have been filtered in validation: %s", data.branchToDeleteType))
	}
	localBranchNameToDelete := data.branchToDeleteInfo.LocalBranchName()
	if _, hasOverride := data.config.NormalConfig.BranchTypeOverrides[localBranchNameToDelete]; hasOverride {
		prog.Value.Add(&opcodes.ConfigRemove{
			Key:   configdomain.NewBranchTypeOverrideKeyForBranch(localBranchNameToDelete).Key,
			Scope: configdomain.ConfigScopeLocal,
		})
	}
	localBranchNameToDelete, hasLocalBranchToDelete := data.branchToDeleteInfo.LocalName.Get()
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         hasLocalBranchToDelete && data.initialBranch != localBranchNameToDelete && data.hasOpenChanges,
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch, Some(data.initialBranch)},
	})
	return optimizer.Optimize(prog.Immutable()), undoProg.Immutable()
}

func deleteFeatureBranch(prog, finalUndoProgram Mutable[program.Program], data deleteData) {
	trackingBranchToDelete, hasTrackingBranchToDelete := data.branchToDeleteInfo.RemoteName.Get()
	if data.branchToDeleteInfo.SyncStatus != gitdomain.SyncStatusDeletedAtRemote && hasTrackingBranchToDelete && data.config.NormalConfig.Offline.IsOnline() {
		ship.UpdateChildBranchProposalsToGrandParent(prog.Value, data.proposalsOfChildBranches)
		prog.Value.Add(&opcodes.BranchTrackingDelete{Branch: trackingBranchToDelete})
	}
	deleteLocalBranch(prog, finalUndoProgram, data)
}

func deleteLocalBranch(prog, finalUndoProgram Mutable[program.Program], data deleteData) {
	if localBranchToDelete, hasLocalBranchToDelete := data.branchToDeleteInfo.LocalName.Get(); hasLocalBranchToDelete {
		if data.initialBranch == localBranchToDelete {
			if data.hasOpenChanges {
				prog.Value.Add(&opcodes.ChangesStage{})
				prog.Value.Add(&opcodes.CommitWithMessage{
					AuthorOverride: None[gitdomain.Author](),
					CommitHook:     configdomain.CommitHookEnabled,
					Message:        "Committing open changes on deleted branch",
				})
				// update the registered initial SHA for this branch so that undo restores the just committed changes
				prog.Value.Add(&opcodes.SnapshotInitialUpdateLocalSHAIfNeeded{Branch: data.initialBranch})
				// when undoing, manually undo the just committed changes so that they are uncommitted again
				finalUndoProgram.Value.Add(&opcodes.CheckoutIfNeeded{Branch: localBranchToDelete})
				finalUndoProgram.Value.Add(&opcodes.UndoLastCommit{})
			}
		}
		// delete the commits of this branch from all descendents
		if data.config.NormalConfig.SyncFeatureStrategy == configdomain.SyncFeatureStrategyRebase {
			descendents := data.config.NormalConfig.Lineage.Descendants(localBranchToDelete)
			for _, descendent := range descendents {
				if branchInfo, hasBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(descendent).Get(); hasBranchInfo {
					sync.RemoveAncestorCommits(sync.RemoveAncestorCommitsArgs{
						Ancestor:          localBranchToDelete.BranchName(),
						Branch:            descendent,
						HasTrackingBranch: branchInfo.HasTrackingBranch(),
						Program:           prog,
						RebaseOnto:        data.config.ValidatedConfigData.MainBranch,
					})
				}
			}
		}
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.branchWhenDone})
		prog.Value.Add(&opcodes.BranchLocalDelete{
			Branch: localBranchToDelete,
		})
		if !data.dryRun {
			sync.RemoveBranchConfiguration(sync.RemoveBranchConfigurationArgs{
				Branch:  localBranchToDelete,
				Lineage: data.config.NormalConfig.Lineage,
				Program: prog,
			})
		}
	}
}

func determineBranchWhenDone(args branchWhenDoneArgs) gitdomain.LocalBranchName {
	if args.branchNameToDelete != args.initialBranch {
		return args.initialBranch
	}
	// here we are deleting the initial branch
	previousBranch, hasPreviousBranch := args.previousBranch.Get()
	if !hasPreviousBranch || previousBranch == args.initialBranch {
		return args.mainBranch
	}
	// here we could return the previous branch
	if previousBranchInfo, hasPreviousBranchInfo := args.branches.FindByLocalName(previousBranch).Get(); hasPreviousBranchInfo {
		if previousBranchInfo.SyncStatus != gitdomain.SyncStatusOtherWorktree {
			return previousBranch
		}
	}
	// here the previous branch is checked out in another worktree --> cannot return it
	return args.mainBranch
}

type branchWhenDoneArgs struct {
	branchNameToDelete gitdomain.LocalBranchName
	branches           gitdomain.BranchInfos
	initialBranch      gitdomain.LocalBranchName
	mainBranch         gitdomain.LocalBranchName
	previousBranch     Option[gitdomain.LocalBranchName]
}

func validateDeleteData(data deleteData) error {
	switch data.branchToDeleteType {
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
	case configdomain.BranchTypeMainBranch:
		return errors.New(messages.DeleteCannotDeleteMainBranch)
	case configdomain.BranchTypePerennialBranch:
		return errors.New(messages.DeleteCannotDeletePerennialBranches)
	}
	return nil
}
