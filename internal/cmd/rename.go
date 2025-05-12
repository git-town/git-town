package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v20/internal/cli/dialog/components"
	"github.com/git-town/git-town/v20/internal/cli/flags"
	"github.com/git-town/git-town/v20/internal/cli/print"
	"github.com/git-town/git-town/v20/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v20/internal/cmd/ship"
	"github.com/git-town/git-town/v20/internal/config"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/execute"
	"github.com/git-town/git-town/v20/internal/forge"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/undo/undoconfig"
	"github.com/git-town/git-town/v20/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v20/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v20/internal/vm/opcodes"
	"github.com/git-town/git-town/v20/internal/vm/optimizer"
	"github.com/git-town/git-town/v20/internal/vm/program"
	"github.com/git-town/git-town/v20/internal/vm/runstate"
	. "github.com/git-town/git-town/v20/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	renameDesc = "Rename a branch and its tracking branch"
	renameHelp = `
The branch to rename must be fully synced.

Renaming perennial branches requires the --force flag.
`
)

func renameCommand() *cobra.Command {
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addForceFlag, readForceFlag := flags.Force("force rename of perennial branch")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "rename [<old_branch_name>] <new_branch_name>",
		Args:  cobra.RangeArgs(1, 2),
		Short: renameDesc,
		Long:  cmdhelpers.Long(renameDesc, renameHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			force, err := readForceFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeRename(args, dryRun, force, verbose)
		},
	}
	addDryRunFlag(&cmd)
	addForceFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRename(args []string, dryRun configdomain.DryRun, force configdomain.Force, verbose configdomain.Verbose) error {
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
	data, exit, err := determineRenameData(args, force, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := renameProgram(data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               "rename",
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

type renameData struct {
	branchInfosLastRun       Option[gitdomain.BranchInfos]
	branchesSnapshot         gitdomain.BranchesSnapshot
	config                   config.ValidatedConfig
	connector                Option[forgedomain.Connector]
	dialogTestInputs         components.TestInputs
	dryRun                   configdomain.DryRun
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	newBranch                gitdomain.LocalBranchName
	nonExistingBranches      gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	oldBranch                gitdomain.BranchInfo
	previousBranch           Option[gitdomain.LocalBranchName]
	proposal                 Option[forgedomain.Proposal]
	proposalsOfChildBranches []forgedomain.Proposal
	stashSize                gitdomain.StashSize
}

func determineRenameData(args []string, force configdomain.Force, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data renameData, exit bool, err error) {
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
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
	connectorOpt, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{oldBranchName},
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
		return data, exit, err
	}
	if validatedConfig.ValidatedConfigData.IsMainBranch(oldBranchName) {
		return data, false, errors.New(messages.RenameMainBranch)
	}
	if !force {
		if validatedConfig.BranchType(oldBranchName) == configdomain.BranchTypePerennialBranch {
			return data, false, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return data, false, errors.New(messages.RenameToSameName)
	}
	if oldBranch.SyncStatus != gitdomain.SyncStatusUpToDate && oldBranch.SyncStatus != gitdomain.SyncStatusLocalOnly {
		return data, false, fmt.Errorf(messages.BranchNotInSyncWithParent, oldBranchName)
	}
	if branchesSnapshot.Branches.HasLocalBranch(newBranchName) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, newBranchName)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(newBranchName, repo.UnvalidatedConfig.NormalConfig.DevRemote) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, newBranchName)
	}
	parentOpt := validatedConfig.NormalConfig.Lineage.Parent(initialBranch)
	lineageBranches := validatedConfig.NormalConfig.Lineage.BranchNames()
	_, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, lineageBranches...)
	proposalOpt := None[forgedomain.Proposal]()
	if !repo.IsOffline {
		proposalOpt = ship.FindProposal(connectorOpt, initialBranch, parentOpt)
	}
	proposalsOfChildBranches := ship.LoadProposalsOfChildBranches(ship.LoadProposalsOfChildBranchesArgs{
		ConnectorOpt:               connectorOpt,
		Lineage:                    validatedConfig.NormalConfig.Lineage,
		Offline:                    false,
		OldBranch:                  oldBranchName,
		OldBranchHasTrackingBranch: oldBranch.HasTrackingBranch(),
	})
	return renameData{
		branchInfosLastRun:       branchInfosLastRun,
		branchesSnapshot:         branchesSnapshot,
		config:                   validatedConfig,
		connector:                connectorOpt,
		dialogTestInputs:         dialogTestInputs,
		dryRun:                   dryRun,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		newBranch:                newBranchName,
		nonExistingBranches:      nonExistingBranches,
		oldBranch:                *oldBranch,
		previousBranch:           previousBranch,
		proposal:                 proposalOpt,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
	}, false, err
}

func renameProgram(data renameData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages)
	oldLocalBranch, hasOldLocalBranch := data.oldBranch.LocalName.Get()
	if !hasOldLocalBranch {
		return prog.Immutable()
	}
	prog.Value.Add(&opcodes.BranchLocalRename{OldName: oldLocalBranch, NewName: data.newBranch})
	if data.initialBranch == oldLocalBranch {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.newBranch})
	}
	if !data.dryRun {
		if override, hasBranchTypeOverride := data.config.NormalConfig.BranchTypeOverrides[oldLocalBranch]; hasBranchTypeOverride {
			prog.Value.Add(
				&opcodes.ConfigSet{
					Key:   configdomain.NewBranchTypeOverrideKeyForBranch(data.newBranch).Key,
					Scope: configdomain.ConfigScopeLocal,
					Value: override.String(),
				},
				&opcodes.ConfigRemove{
					Key:   configdomain.NewBranchTypeOverrideKeyForBranch(oldLocalBranch).Key,
					Scope: configdomain.ConfigScopeLocal,
				},
			)
		}
		if parentBranch, hasParent := data.config.NormalConfig.Lineage.Parent(oldLocalBranch).Get(); hasParent {
			prog.Value.Add(&opcodes.LineageParentSet{Branch: data.newBranch, Parent: parentBranch})
		}
		prog.Value.Add(&opcodes.LineageParentRemove{Branch: oldLocalBranch})
	}
	for _, child := range data.config.NormalConfig.Lineage.Children(oldLocalBranch) {
		prog.Value.Add(&opcodes.LineageParentSet{Branch: child, Parent: data.newBranch})
	}
	if oldTrackingBranch, hasOldTrackingBranch := data.oldBranch.RemoteName.Get(); hasOldTrackingBranch {
		if data.oldBranch.HasTrackingBranch() && data.config.NormalConfig.Offline.IsOnline() {
			prog.Value.Add(&opcodes.BranchTrackingCreate{Branch: data.newBranch})
			updateChildBranchProposalsToBranch(prog.Value, data.proposalsOfChildBranches, data.newBranch)
			proposal, hasProposal := data.proposal.Get()
			connector, hasConnector := data.connector.Get()
			connectorCanUpdateProposalSource := hasConnector && connector.UpdateProposalSourceFn().IsSome()
			if hasProposal && hasConnector && connectorCanUpdateProposalSource {
				prog.Value.Add(&opcodes.ProposalUpdateSource{
					NewBranch:      data.newBranch,
					OldBranch:      data.oldBranch.LocalBranchName(),
					ProposalNumber: proposal.Number,
				})
			}
			prog.Value.Add(&opcodes.BranchTrackingDelete{Branch: oldTrackingBranch})
		}
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{Some(data.newBranch), data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             false,
		StashOpenChanges:         false,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return optimizer.Optimize(prog.Immutable())
}

func updateChildBranchProposalsToBranch(prog *program.Program, proposals []forgedomain.Proposal, target gitdomain.LocalBranchName) {
	for _, childProposal := range proposals {
		prog.Add(&opcodes.ProposalUpdateTarget{
			NewBranch:      target,
			OldBranch:      childProposal.Target,
			ProposalNumber: childProposal.Number,
		})
	}
}
