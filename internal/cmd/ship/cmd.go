package ship

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/cmd/sync"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/validate"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/optimizer"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addIgnoreUncommittedFlag, readIgnoreUncommittedFlag := flags.IgnoreUncommitted()
	addMessageFileFlag, readMessageFileFlag := flags.CommitMessageFile()
	addMessageFlag, readMessageFlag := flags.CommitMessage("specify the commit message for the squash commit")
	addShipStrategyFlag, readShipStrategyFlag := flags.ShipStrategy()
	addToParentFlag, readToParentFlag := flags.ShipIntoNonPerennialParent()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   shipCommand,
		Args:  cobra.MaximumNArgs(1),
		Short: shipDesc,
		Long:  cmdhelpers.Long(shipDesc, fmt.Sprintf(shipHelp, configdomain.KeyGitHubToken)),
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, errDryRun := readDryRunFlag(cmd)
			ignoreUncommitted, errIgnoreUncommitted := readIgnoreUncommittedFlag(cmd)
			message, errMessage := readMessageFlag(cmd)
			messageFile, errMessageFile := readMessageFileFlag(cmd)
			shipStrategy, errShipStrategy := readShipStrategyFlag(cmd)
			toParent, errToParent := readToParentFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errDryRun, errIgnoreUncommitted, errMessage, errMessageFile, errShipStrategy, errToParent, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            dryRun,
				IgnoreUncommitted: ignoreUncommitted,
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeShip(executeShipArgs{
				args:         args,
				cliConfig:    cliConfig,
				message:      message,
				messageFile:  messageFile,
				shipStrategy: shipStrategy,
				toParent:     toParent,
			})
		},
	}
	addMessageFileFlag(&cmd)
	addDryRunFlag(&cmd)
	addIgnoreUncommittedFlag(&cmd)
	addMessageFlag(&cmd)
	addShipStrategyFlag(&cmd)
	addToParentFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type executeShipArgs struct {
	args         []string
	cliConfig    configdomain.PartialConfig
	message      Option[gitdomain.CommitMessage]
	messageFile  Option[gitdomain.CommitMessageFile]
	shipStrategy Option[configdomain.ShipStrategy]
	toParent     configdomain.ShipIntoNonperennialParent
}

func executeShip(args executeShipArgs) error {
Start:
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	sharedData, flow, err := determineSharedShipData(determineSharedShipDataArgs{
		args:                 args.args,
		repo:                 repo,
		shipStrategyOverride: args.shipStrategy,
	})
	if err != nil {
		return err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit:
		return nil
	case configdomain.ProgramFlowRestart:
		goto Start
	}
	message, err := ReadFile(args.message, args.messageFile)
	if err != nil {
		return err
	}
	if err = validateSharedData(sharedData, args.toParent, message); err != nil {
		return err
	}
	prog := NewMutable(&program.Program{})
	switch sharedData.config.NormalConfig.ShipStrategy {
	case configdomain.ShipStrategyAPI:
		apiData, err := determineAPIData(sharedData)
		if err != nil {
			return err
		}
		if err = shipAPIProgram(prog, repo, sharedData, apiData, message); err != nil {
			return err
		}
	case configdomain.ShipStrategyAlwaysMerge:
		mergeData, err := determineMergeData(repo, sharedData.branchToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramAlwaysMerge(repo, shipProgramAlwaysMergeArgs{
			commitMessage: message,
			mergeData:     mergeData,
			prog:          prog,
			sharedData:    sharedData,
		})
	case configdomain.ShipStrategyFastForward:
		mergeData, err := determineMergeData(repo, sharedData.branchToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramFastForward(prog, repo, sharedData, mergeData)
	case configdomain.ShipStrategySquashMerge:
		squashMergeData, err := determineMergeData(repo, sharedData.branchToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramSquashMerge(prog, repo, sharedData, squashMergeData, message)
	}
	if sharedData.config.NormalConfig.ProposalsShowLineage == forgedomain.ProposalsShowLineageCLI {
		_ = sync.AddStackLineageUpdateOpcodes(
			sync.AddStackLineageUpdateOpcodesArgs{
				Current:   sharedData.initialBranch,
				FullStack: true,
				Program:   prog,
				ProposalStackLineageArgs: proposallineage.ProposalStackLineageArgs{
					Connector:                forgedomain.ProposalFinderFromConnector(sharedData.connector),
					CurrentBranch:            sharedData.initialBranch,
					Lineage:                  sharedData.config.NormalConfig.Lineage,
					MainAndPerennialBranches: sharedData.config.MainAndPerennials(),
					Order:                    sharedData.config.NormalConfig.Order,
				},
				ProposalStackLineageTree: None[*proposallineage.Tree](),
				// Proposal has been shipped and its stack lineage
				// information shouldn't need to be updated because
				// proposal is not in a review state.
				SkipUpdateForProposalsWithBaseBranch: gitdomain.LocalBranchNames{sharedData.initialBranch},
			},
		)
	}
	// Stash uncommitted changes if ignore-uncommitted is enabled
	shouldStash := sharedData.hasOpenChanges && sharedData.config.NormalConfig.ShipIgnoreUncommitted.AllowUncommitted()
	if shouldStash {
		cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
			DryRun:                   sharedData.config.NormalConfig.DryRun,
			InitialStashSize:         sharedData.stashSize,
			PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{sharedData.previousBranch},
			RunInGitRoot:             false,
			StashOpenChanges:         shouldStash,
		})
	}
	optimizedProgram := optimizer.Optimize(prog.Immutable())
	runState := runstate.RunState{
		BeginBranchesSnapshot: sharedData.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        sharedData.stashSize,
		Command:               shipCommand,
		DryRun:                sharedData.config.NormalConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[configdomain.EndConfigSnapshot](),
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
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          sharedData.hasOpenChanges,
		InitialBranch:           sharedData.initialBranch,
		InitialBranchesSnapshot: sharedData.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        sharedData.stashSize,
		Inputs:                  sharedData.inputs,
		PendingCommand:          None[string](),
		RootDir:                 repo.RootDir,
		RunState:                runState,
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
		branch := data.branchToShipInfo.LocalName.GetOrPanic()
		parentBranch := data.targetBranch.LocalName.GetOrPanic()
		if !data.config.IsMainOrPerennialBranch(parentBranch) {
			ancestors := data.config.NormalConfig.Lineage.Ancestors(branch)
			ancestorsWithoutMainOrPerennial := ancestors[1:]
			oldestAncestor := ancestorsWithoutMainOrPerennial[0]
			return fmt.Errorf(messages.ShipChildBranch, stringslice.Connect(ancestorsWithoutMainOrPerennial.Strings()), oldestAncestor)
		}
	}
	switch data.branchToShipInfo.SyncStatus {
	case gitdomain.SyncStatusDeletedAtRemote:
		return fmt.Errorf(messages.BranchDeletedAtRemote, data.branchToShip)
	case
		gitdomain.SyncStatusNotInSync,
		gitdomain.SyncStatusAhead,
		gitdomain.SyncStatusBehind:
		return fmt.Errorf(messages.BranchNotInSyncWithParent, data.branchToShip)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.ShipBranchIsInOtherWorktree, data.branchToShip)
	case
		gitdomain.SyncStatusUpToDate,
		gitdomain.SyncStatusRemoteOnly,
		gitdomain.SyncStatusLocalOnly:
	}
	if localName, hasLocalName := data.branchToShipInfo.LocalName.Get(); hasLocalName {
		if localName == data.initialBranch {
			if !data.config.NormalConfig.ShipIgnoreUncommitted.AllowUncommitted() {
				return validate.NoOpenChanges(data.hasOpenChanges)
			}
		}
	}
	return nil
}

func ReadFile[TEXT ~string, FILE FileFlag](inputText Option[TEXT], inputFileOpt Option[FILE]) (Option[TEXT], error) {
	if inputText.IsSome() {
		return inputText, nil
	}
	file, hasFile := inputFileOpt.Get()
	if !hasFile {
		return None[TEXT](), nil
	}
	if file.ShouldReadStdin() {
		content, err := io.ReadAll(os.Stdin)
		return NewOption(TEXT(string(content))), gohacks.WrapIfError(err, "cannot read STDIN: %w")
	}
	fileData, err := os.ReadFile(file.String())
	return NewOption(TEXT(string(fileData))), err
}

type FileFlag interface {
	ShouldReadStdin() bool
	fmt.Stringer
}
