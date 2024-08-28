package ship

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/flags"
	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/execute"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/gohacks/stringslice"
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

const shipCommand = "ship"

const shipDesc = "Deliver a completed feature branch"

const shipHelp = `
Merges the given or current feature branch into its parent.
How exactly this happen depends on the configured ship-strategy.

Ships only direct children of the main branch.
To ship a child branch, ship or kill all ancestor branches first
or ship with the "--to-parent" flag.

To use the online functionality, configure a personal access token with the "repo" scope
and run Run 'git config %s <token>' (optionally add the '--global' flag).

If your origin server deletes shipped branches,
disable the ship-delete-tracking-branch configuration setting.`

func Cmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addMessageFlag, readMessageFlag := flags.CommitMessage("specify the commit message for the squash commit")
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addToParentFlag, readToParentFlag := flags.ShipIntoNonPerennialParent()
	cmd := cobra.Command{
		Use:   shipCommand,
		Args:  cobra.MaximumNArgs(1),
		Short: shipDesc,
		Long:  cmdhelpers.Long(shipDesc, fmt.Sprintf(shipHelp, configdomain.KeyGithubToken)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeShip(args, readMessageFlag(cmd), readDryRunFlag(cmd), readVerboseFlag(cmd), readToParentFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addMessageFlag(&cmd)
	addToParentFlag(&cmd)
	return &cmd
}

func executeShip(args []string, message Option[gitdomain.CommitMessage], dryRun configdomain.DryRun, verbose configdomain.Verbose, toParent configdomain.ShipIntoNonperennialParent) error {
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
	sharedData, exit, err := determineSharedShipData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateSharedData(sharedData, toParent)
	if err != nil {
		return err
	}
	var shipProgram program.Program
	switch sharedData.config.Config.ShipStrategy {
	case configdomain.ShipStrategyAPI:
		apiData, err := determineAPIData(sharedData)
		if err != nil {
			return err
		}
		shipProgram = shipAPIProgram(sharedData, apiData, message)
	case configdomain.ShipStrategySquashMerge:
		squashMergeData, err := determineSquashMergeData(repo)
		if err != nil {
			return err
		}
		shipProgram = shipProgramSquashMerge(sharedData, squashMergeData, message)
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: sharedData.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        sharedData.stashSize,
		Command:               shipCommand,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            shipProgram,
		TouchedBranches:       shipProgram.TouchedBranches(),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  sharedData.config,
		Connector:               sharedData.connector,
		DialogTestInputs:        sharedData.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          sharedData.hasOpenChanges,
		InitialBranch:           sharedData.initialBranch,
		InitialBranchesSnapshot: sharedData.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        sharedData.stashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

func ensureParentBranchIsMainOrPerennialBranch(branch, parentBranch gitdomain.LocalBranchName, config configdomain.ValidatedConfig, lineage configdomain.Lineage) error {
	if !config.IsMainOrPerennialBranch(parentBranch) {
		ancestors := lineage.Ancestors(branch)
		ancestorsWithoutMainOrPerennial := ancestors[1:]
		oldestAncestor := ancestorsWithoutMainOrPerennial[0]
		return fmt.Errorf(messages.ShipChildBranch, stringslice.Connect(ancestorsWithoutMainOrPerennial.Strings()), oldestAncestor)
	}
	return nil
}

func updateChildBranchProposals(prog *program.Program, proposals []hostingdomain.Proposal, targetBranch gitdomain.LocalBranchName) {
	for _, childProposal := range proposals {
		prog.Add(&opcodes.UpdateProposalTarget{
			ProposalNumber: childProposal.Number,
			NewTarget:      targetBranch,
		})
	}
}

func validateShippableBranchType(branchType configdomain.BranchType) error {
	switch branchType {
	case configdomain.BranchTypeContributionBranch:
		return errors.New(messages.ContributionBranchCannotShip)
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
		return nil
	case configdomain.BranchTypeMainBranch:
		return errors.New(messages.MainBranchCannotShip)
	case configdomain.BranchTypeObservedBranch:
		return errors.New(messages.ObservedBranchCannotShip)
	case configdomain.BranchTypePerennialBranch:
		return errors.New(messages.PerennialBranchCannotShip)
	}
	panic(fmt.Sprintf("unhandled branch type: %v", branchType))
}

func validateSyncStatus(syncStatus gitdomain.SyncStatus, branchToShip gitdomain.LocalBranchName) error {
	switch syncStatus {
	case gitdomain.SyncStatusDeletedAtRemote:
		return fmt.Errorf(messages.ShipBranchDeletedAtRemote, branchToShip)
	case gitdomain.SyncStatusNotInSync:
		return fmt.Errorf(messages.ShipBranchNotInSync, branchToShip)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.ShipBranchIsInOtherWorktree, branchToShip)
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusRemoteOnly, gitdomain.SyncStatusLocalOnly:
	}
	return nil
}

func validateSharedData(data sharedShipData, toParent configdomain.ShipIntoNonperennialParent) error {
	if !toParent {
		err := ensureParentBranchIsMainOrPerennialBranch(data.branchToShip.LocalName.GetOrPanic(), data.targetBranch.LocalName.GetOrPanic(), data.config.Config, data.config.Config.Lineage)
		if err != nil {
			return err
		}
	}
	if err := validateSyncStatus(data.branchToShip.SyncStatus, data.branchNameToShip); err != nil {
		return err
	}
	if localName, hasLocalName := data.branchToShip.LocalName.Get(); hasLocalName {
		if localName == data.initialBranch {
			return validate.NoOpenChanges(data.hasOpenChanges)
		}
	}
	if err := validateSyncStatus(data.branchToShip.SyncStatus, data.branchNameToShip); err != nil {
		return err
	}
	return nil
}
