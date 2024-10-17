package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/cmd/ship"
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

const renameDesc = "Rename a branch both locally and remotely"

const renameHelp = `
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
			return executeRename(args, readDryRunFlag(cmd), readForceFlag(cmd), readVerboseFlag(cmd))
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
	runProgram := renameProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
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
	branchesSnapshot         gitdomain.BranchesSnapshot
	config                   config.ValidatedConfig
	connector                Option[hostingdomain.Connector]
	dialogTestInputs         components.TestInputs
	dryRun                   configdomain.DryRun
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	newBranch                gitdomain.LocalBranchName
	oldBranch                gitdomain.BranchInfo
	previousBranch           Option[gitdomain.LocalBranchName]
	proposal                 Option[hostingdomain.Proposal]
	proposalsOfChildBranches []hostingdomain.Proposal
	stashSize                gitdomain.StashSize
}

func determineRenameData(args []string, force configdomain.Force, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data renameData, exit bool, err error) {
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
	connectorOpt, err := hosting.NewConnector(repo.UnvalidatedConfig, gitdomain.RemoteOrigin, print.Logger{})
	if err != nil {
		return data, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
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
		return data, false, fmt.Errorf(messages.RenameNotInSync, oldBranchName)
	}
	if branchesSnapshot.Branches.HasLocalBranch(newBranchName) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, newBranchName)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(newBranchName) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, newBranchName)
	}
	parentOpt := validatedConfig.Config.Lineage.Parent(initialBranch)
	proposalOpt := ship.FindProposal(connectorOpt, initialBranch, parentOpt)
	proposalsOfChildBranches := ship.LoadProposalsOfChildBranches(ship.LoadProposalsOfChildBranchesArgs{
		ConnectorOpt:               connectorOpt,
		Lineage:                    validatedConfig.Config.Lineage,
		Offline:                    false,
		OldBranch:                  oldBranchName,
		OldBranchHasTrackingBranch: oldBranch.HasTrackingBranch(),
	})
	return renameData{
		branchesSnapshot:         branchesSnapshot,
		config:                   validatedConfig,
		connector:                connectorOpt,
		dialogTestInputs:         dialogTestInputs,
		dryRun:                   dryRun,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		newBranch:                newBranchName,
		oldBranch:                *oldBranch,
		previousBranch:           previousBranch,
		proposal:                 proposalOpt,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
	}, false, err
}

func renameProgram(data renameData) program.Program {
	result := NewMutable(&program.Program{})
	oldLocalBranch, hasOldLocalBranch := data.oldBranch.LocalName.Get()
	if !hasOldLocalBranch {
		return result.Get()
	}
	result.Value.Add(&opcodes.BranchLocalRename{OldName: oldLocalBranch, NewName: data.newBranch})
	if data.initialBranch == oldLocalBranch {
		result.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.newBranch})
	}
	if !data.dryRun {
		if data.config.Config.IsPerennialBranch(data.initialBranch) {
			result.Value.Add(&opcodes.BranchesPerennialRemove{Branch: oldLocalBranch})
			result.Value.Add(&opcodes.BranchesPerennialAdd{Branch: data.newBranch})
		} else {
			if slices.Contains(data.config.Config.PrototypeBranches, data.initialBranch) {
				result.Value.Add(&opcodes.BranchesPrototypeRemove{Branch: oldLocalBranch})
				result.Value.Add(&opcodes.BranchesPrototypeAdd{Branch: data.newBranch})
			}
			if slices.Contains(data.config.Config.ObservedBranches, data.initialBranch) {
				result.Value.Add(&opcodes.BranchesObservedRemove{Branch: oldLocalBranch})
				result.Value.Add(&opcodes.BranchesObservedAdd{Branch: data.newBranch})
			}
			if slices.Contains(data.config.Config.ContributionBranches, data.initialBranch) {
				result.Value.Add(&opcodes.BranchesContributionRemove{Branch: oldLocalBranch})
				result.Value.Add(&opcodes.BranchesContributionAdd{Branch: data.newBranch})
			}
			if slices.Contains(data.config.Config.ParkedBranches, data.initialBranch) {
				result.Value.Add(&opcodes.BranchesParkedRemove{Branch: oldLocalBranch})
				result.Value.Add(&opcodes.BranchesParkedAdd{Branch: data.newBranch})
			}
			if parentBranch, hasParent := data.config.Config.Lineage.Parent(oldLocalBranch).Get(); hasParent {
				result.Value.Add(&opcodes.LineageParentSet{Branch: data.newBranch, Parent: parentBranch})
			}
			result.Value.Add(&opcodes.BranchParentDelete{Branch: oldLocalBranch})
		}
	}
	for _, child := range data.config.Config.Lineage.Children(oldLocalBranch) {
		result.Value.Add(&opcodes.LineageParentSet{Branch: child, Parent: data.newBranch})
	}
	if oldTrackingBranch, hasOldTrackingBranch := data.oldBranch.RemoteName.Get(); hasOldTrackingBranch {
		if data.oldBranch.HasTrackingBranch() && data.config.Config.IsOnline() {
			result.Value.Add(&opcodes.BranchTrackingCreate{Branch: data.newBranch})
			updateChildBranchProposalsToBranch(result.Value, data.proposalsOfChildBranches, data.newBranch)
			if proposal, hasProposal := data.proposal.Get(); hasProposal {
				result.Value.Add(&opcodes.ProposalUpdateHead{
					NewTarget:      data.newBranch,
					OldTarget:      data.oldBranch.LocalBranchName(),
					ProposalNumber: proposal.Number,
				})
			}
			result.Value.Add(&opcodes.BranchTrackingDelete{Branch: oldTrackingBranch})
		}
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{Some(data.newBranch), data.previousBranch}
	cmdhelpers.Wrap(result, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             false,
		StashOpenChanges:         false,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return result.Get()
}

func updateChildBranchProposalsToBranch(prog *program.Program, proposals []hostingdomain.Proposal, target gitdomain.LocalBranchName) {
	for _, childProposal := range proposals {
		prog.Add(&opcodes.ProposalUpdateBase{
			NewTarget:      target,
			OldTarget:      childProposal.Target,
			ProposalNumber: childProposal.Number,
		})
	}
}
