package ship

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
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
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addMessageFlag, readMessageFlag := flags.CommitMessage("specify the commit message for the squash commit")
	addMessageFileFlag, readMessageFileFlag := flags.CommitMessageFile()
	addShipStrategyFlag, readShipStrategyFlag := flags.ShipStrategy()
	addToParentFlag, readToParentFlag := flags.ShipIntoNonPerennialParent()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   shipCommand,
		Args:  cobra.MaximumNArgs(1),
		Short: shipDesc,
		Long:  cmdhelpers.Long(shipDesc, fmt.Sprintf(shipHelp, configdomain.KeyGitHubToken)),
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, err1 := readDryRunFlag(cmd)
			message, err2 := readMessageFlag(cmd)
			messageFile, err3 := readMessageFileFlag(cmd)
			shipStrategy, err4 := readShipStrategyFlag(cmd)
			toParent, err5 := readToParentFlag(cmd)
			verbose, err6 := readVerboseFlag(cmd)
			if err := cmp.Or(err1, err2, err3, err4, err5, err6); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				DryRun:  dryRun,
				Verbose: verbose,
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
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	sharedData, exit, err := determineSharedShipData(args.args, repo, args.shipStrategy)
	if err != nil || exit {
		return err
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
		mergeData, err := determineMergeData(repo, sharedData.branchNameToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramAlwaysMerge(shipProgramAlwaysMergeArgs{
			prog:          prog,
			repo:          repo,
			sharedData:    sharedData,
			mergeData:     mergeData,
			commitMessage: message,
		})
	case configdomain.ShipStrategyFastForward:
		mergeData, err := determineMergeData(repo, sharedData.branchNameToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramFastForward(prog, repo, sharedData, mergeData)
	case configdomain.ShipStrategySquashMerge:
		squashMergeData, err := determineMergeData(repo, sharedData.branchNameToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramSquashMerge(prog, repo, sharedData, squashMergeData, message)
	}
	optimizedProgram := optimizer.Optimize(prog.Immutable())
	runState := runstate.RunState{
		BeginBranchesSnapshot: sharedData.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        sharedData.stashSize,
		Command:               shipCommand,
		DryRun:                sharedData.config.NormalConfig.DryRun,
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
