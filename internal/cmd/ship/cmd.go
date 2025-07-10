package ship

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
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
	shipCommand = "ship"
	shipDesc    = "Deliver a completed feature branch"
	shipHelp    = `
Merges the given or current feature branch into its parent.
How exactly this happen depends on the configured ship-strategy.

Ships only direct children of the main branch.
To ship a child branch, ship or delete all ancestor branches first
or ship with the "--to-parent" flag.

To use the online functionality,
configure a personal access token with the "repo" scope
and run "git config %s <token>" (optionally add the "--global" flag).

If your origin server deletes shipped branches,
disable the ship-delete-tracking-branch configuration setting.`
)

func Cmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addMessageFlag, readMessageFlag := flags.CommitMessage("specify the commit message for the squash commit")
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addShipStrategyFlag, readShipStrategyFlag := flags.ShipStrategy()
	addToParentFlag, readToParentFlag := flags.ShipIntoNonPerennialParent()
	cmd := cobra.Command{
		Use:   shipCommand,
		Args:  cobra.MaximumNArgs(1),
		Short: shipDesc,
		Long:  cmdhelpers.Long(shipDesc, fmt.Sprintf(shipHelp, configdomain.KeyGitHubToken)),
		RunE: func(cmd *cobra.Command, args []string) error {
			shipStrategyOverride, err1 := readShipStrategyFlag(cmd)
			message, err2 := readMessageFlag(cmd)
			dryRun, err3 := readDryRunFlag(cmd)
			toParent, err4 := readToParentFlag(cmd)
			verbose, err5 := readVerboseFlag(cmd)
			if err := cmp.Or(err1, err2, err3, err4, err5); err != nil {
				return err
			}
			cliConfig := cliconfig.CliConfig{
				DryRun:  dryRun,
				Verbose: verbose,
			}
			return executeShip(args, cliConfig, message, shipStrategyOverride, toParent)
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addMessageFlag(&cmd)
	addShipStrategyFlag(&cmd)
	addToParentFlag(&cmd)
	return &cmd
}

func executeShip(args []string, cliConfig cliconfig.CliConfig, message Option[gitdomain.CommitMessage], shipStrategy Option[configdomain.ShipStrategy], toParent configdomain.ShipIntoNonperennialParent) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	sharedData, exit, err := determineSharedShipData(args, repo, cliConfig, shipStrategy)
	if err != nil || exit {
		return err
	}
	if err = validateSharedData(sharedData, toParent, message); err != nil {
		return err
	}
	prog := NewMutable(&program.Program{})
	switch sharedData.config.NormalConfig.ShipStrategy {
	case configdomain.ShipStrategyAPI:
		apiData, err := determineAPIData(sharedData)
		if err != nil {
			return err
		}
		if err = shipAPIProgram(prog, sharedData, apiData, message); err != nil {
			return err
		}
	case configdomain.ShipStrategyAlwaysMerge:
		mergeData, err := determineMergeData(repo, sharedData.branchNameToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramAlwaysMerge(prog, sharedData, mergeData, message)
	case configdomain.ShipStrategyFastForward:
		mergeData, err := determineMergeData(repo, sharedData.branchNameToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramFastForward(prog, sharedData, mergeData)
	case configdomain.ShipStrategySquashMerge:
		squashMergeData, err := determineMergeData(repo, sharedData.branchNameToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramSquashMerge(prog, sharedData, squashMergeData, message)
	}
	optimizedProgram := optimizer.Optimize(prog.Immutable())
	runState := runstate.RunState{
		BeginBranchesSnapshot: sharedData.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        sharedData.stashSize,
		Command:               shipCommand,
		DryRun:                cliConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		BranchInfosLastRun:    sharedData.previousBranchInfos,
		RunProgram:            optimizedProgram,
		TouchedBranches:       optimizedProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  sharedData.config,
		Connector:               sharedData.connector,
		Detached:                false,
		DialogTestInputs:        sharedData.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          sharedData.hasOpenChanges,
		InitialBranch:           sharedData.initialBranch,
		InitialBranchesSnapshot: sharedData.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        sharedData.stashSize,
		PendingCommand:          None[string](),
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 cliConfig.Verbose,
	})
}

func UpdateChildBranchProposalsToGrandParent(prog *program.Program, proposals []forgedomain.Proposal) {
	for _, childProposal := range proposals {
		data := childProposal.Data.Data()
		prog.Add(&opcodes.ProposalUpdateTargetToGrandParent{
			Branch:    data.Source,
			OldTarget: data.Target,
			Proposal:  childProposal,
		})
	}
}

func validateSharedData(data sharedShipData, toParent configdomain.ShipIntoNonperennialParent, message Option[gitdomain.CommitMessage]) error {
	if data.config.NormalConfig.ShipStrategy == configdomain.ShipStrategyFastForward && message.IsSome() {
		return errors.New(messages.ShipMessageWithFastForward)
	}
	if !toParent {
		branch := data.branchToShip.LocalName.GetOrPanic()
		parentBranch := data.targetBranch.LocalName.GetOrPanic()
		if !data.config.IsMainOrPerennialBranch(parentBranch) {
			ancestors := data.config.NormalConfig.Lineage.Ancestors(branch)
			ancestorsWithoutMainOrPerennial := ancestors[1:]
			oldestAncestor := ancestorsWithoutMainOrPerennial[0]
			return fmt.Errorf(messages.ShipChildBranch, stringslice.Connect(ancestorsWithoutMainOrPerennial.Strings()), oldestAncestor)
		}
	}
	switch data.branchToShip.SyncStatus {
	case gitdomain.SyncStatusDeletedAtRemote:
		return fmt.Errorf(messages.BranchDeletedAtRemote, data.branchNameToShip)
	case
		gitdomain.SyncStatusNotInSync,
		gitdomain.SyncStatusAhead,
		gitdomain.SyncStatusBehind:
		return fmt.Errorf(messages.BranchNotInSyncWithParent, data.branchNameToShip)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.ShipBranchIsInOtherWorktree, data.branchNameToShip)
	case
		gitdomain.SyncStatusUpToDate,
		gitdomain.SyncStatusRemoteOnly,
		gitdomain.SyncStatusLocalOnly:
	}
	if localName, hasLocalName := data.branchToShip.LocalName.Get(); hasLocalName {
		if localName == data.initialBranch {
			return validate.NoOpenChanges(data.hasOpenChanges)
		}
	}
	return nil
}
