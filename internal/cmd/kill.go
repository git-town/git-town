package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/cmd/ship"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/slice"
	"github.com/git-town/git-town/v16/internal/hosting"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/sync"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	"github.com/git-town/git-town/v16/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const killDesc = "Remove an obsolete feature branch"

const killHelp = `
Deletes the current or provided branch from the local and origin repositories. Does not delete perennial branches nor the main branch.`

func killCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:   "kill [<branch>]",
		Args:  cobra.MaximumNArgs(1),
		Short: killDesc,
		Long:  cmdhelpers.Long(killDesc, killHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeKill(args, readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeKill(args []string, dryRun configdomain.DryRun, verbose configdomain.Verbose) error {
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
	data, exit, err := determineKillData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateKillData(data)
	if err != nil {
		return err
	}
	runProgram, finalUndoProgram := killProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "kill",
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		FinalUndoProgram:      finalUndoProgram,
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
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

type killData struct {
	branchToKillInfo         gitdomain.BranchInfo
	branchToKillType         configdomain.BranchType
	branchWhenDone           gitdomain.LocalBranchName
	branchesSnapshot         gitdomain.BranchesSnapshot
	config                   config.ValidatedConfig
	connector                Option[hostingdomain.Connector]
	dialogTestInputs         components.TestInputs
	dryRun                   configdomain.DryRun
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	parentBranch             gitdomain.LocalBranchName
	previousBranch           Option[gitdomain.LocalBranchName]
	proposalsOfChildBranches []hostingdomain.Proposal
	stashSize                gitdomain.StashSize
}

func determineKillData(args []string, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data killData, exit bool, err error) {
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
	branchNameToKill := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToKill, hasBranchToKill := branchesSnapshot.Branches.FindByLocalName(branchNameToKill).Get()
	if !hasBranchToKill {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToKill)
	}
	if branchToKill.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return data, exit, fmt.Errorf(messages.KillBranchOtherWorktree, branchNameToKill)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.BranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{},
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
	branchTypeToKill := validatedConfig.Config.BranchType(branchNameToKill)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	previousBranchOpt := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	branchWhenDone := determineBranchWhenDone(branchWhenDoneArgs{
		branchNameToKill: branchNameToKill,
		branches:         branchesSnapshot.Branches,
		initialBranch:    initialBranch,
		mainBranch:       validatedConfig.Config.MainBranch,
		previousBranch:   previousBranchOpt,
	})
	localBranchToKill, hasLocalBranchToKill := branchToKill.LocalName.Get()
	var parentBranch Option[gitdomain.LocalBranchName]
	if hasLocalBranchToKill {
		parentBranch = validatedConfig.Config.Lineage.Parent(localBranchToKill)
	} else {
		parentBranch = None[gitdomain.LocalBranchName]()
	}
	var connectorOpt Option[hostingdomain.Connector]
	if originURL, hasOriginURL := validatedConfig.OriginURL().Get(); hasOriginURL {
		connectorOpt, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          *validatedConfig.Config.UnvalidatedConfig,
			HostingPlatform: validatedConfig.Config.HostingPlatform,
			Log:             print.Logger{},
			RemoteURL:       originURL,
		})
		if err != nil {
			return data, false, err
		}
	}
	var proposalsOfChildBranches []hostingdomain.Proposal
	childBranches := validatedConfig.Config.Lineage.Children(branchNameToKill)
	connector, hasConnector := connectorOpt.Get()
	if hasConnector && connector.CanMakeAPICalls() {
		if !repo.IsOffline.IsTrue() && branchToKill.HasTrackingBranch() {
			for _, childBranch := range childBranches {
				childProposalOpt, err := connector.FindProposal(childBranch, branchNameToKill)
				if err != nil {
					return data, false, fmt.Errorf(messages.ProposalNotFoundForBranch, branchNameToKill, err)
				}
				childProposal, hasChildProposal := childProposalOpt.Get()
				if hasChildProposal {
					proposalsOfChildBranches = append(proposalsOfChildBranches, childProposal)
				}
			}
		}
	}
	mainBranch := validatedConfig.Config.MainBranch
	return killData{
		branchToKillInfo:         *branchToKill,
		branchToKillType:         branchTypeToKill,
		branchWhenDone:           branchWhenDone,
		branchesSnapshot:         branchesSnapshot,
		config:                   validatedConfig,
		connector:                connectorOpt,
		dialogTestInputs:         dialogTestInputs,
		dryRun:                   dryRun,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		parentBranch:             parentBranch.GetOrElse(mainBranch),
		previousBranch:           previousBranchOpt,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
	}, false, nil
}

func killProgram(data killData) (runProgram, finalUndoProgram program.Program) {
	prog := NewMutable(&program.Program{})
	undoProg := NewMutable(&program.Program{})
	switch data.branchToKillType {
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
		killFeatureBranch(prog, undoProg, data)
	case configdomain.BranchTypeObservedBranch, configdomain.BranchTypeContributionBranch:
		killLocalBranch(prog, undoProg, data)
	case configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		panic(fmt.Sprintf("this branch type should have been filtered in validation: %s", data.branchToKillType))
	}
	localBranchNameToKill, hasLocalBranchToKill := data.branchToKillInfo.LocalName.Get()
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         hasLocalBranchToKill && data.initialBranch != localBranchNameToKill && data.hasOpenChanges,
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch, Some(data.initialBranch)},
	})
	return prog.Get(), undoProg.Get()
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killFeatureBranch(prog, finalUndoProgram Mutable[program.Program], data killData) {
	trackingBranchToKill, hasTrackingBranchToKill := data.branchToKillInfo.RemoteName.Get()
	if data.branchToKillInfo.SyncStatus != gitdomain.SyncStatusDeletedAtRemote && hasTrackingBranchToKill && data.config.Config.IsOnline() {
		ship.UpdateChildBranchProposals(prog.Value, data.proposalsOfChildBranches, data.parentBranch)
		prog.Value.Add(&opcodes.DeleteTrackingBranch{Branch: trackingBranchToKill})
	}
	killLocalBranch(prog, finalUndoProgram, data)
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killLocalBranch(prog, finalUndoProgram Mutable[program.Program], data killData) {
	if localBranchToKill, hasLocalBranchToKill := data.branchToKillInfo.LocalName.Get(); hasLocalBranchToKill {
		if data.initialBranch == localBranchToKill {
			if data.hasOpenChanges {
				prog.Value.Add(&opcodes.StageOpenChanges{})
				prog.Value.Add(&opcodes.CommitOpenChanges{Message: "Committing WIP for git town undo"})
				// update the registered initial SHA for this branch so that undo restores the just committed changes
				prog.Value.Add(&opcodes.UpdateInitialBranchLocalSHA{Branch: data.initialBranch})
				// when undoing, manually undo the just committed changes so that they are uncommitted again
				finalUndoProgram.Value.Add(&opcodes.Checkout{Branch: localBranchToKill})
				finalUndoProgram.Value.Add(&opcodes.UndoLastCommit{})
			}
			prog.Value.Add(&opcodes.Checkout{Branch: data.branchWhenDone})
		}
		prog.Value.Add(&opcodes.DeleteLocalBranch{Branch: localBranchToKill})
		if data.dryRun.IsFalse() {
			sync.RemoveBranchFromLineage(sync.RemoveBranchFromLineageArgs{
				Branch:  localBranchToKill,
				Lineage: data.config.Config.Lineage,
				Parent:  data.parentBranch,
				Program: prog,
			})
		}
	}
}

func determineBranchWhenDone(args branchWhenDoneArgs) gitdomain.LocalBranchName {
	if args.branchNameToKill != args.initialBranch {
		return args.initialBranch
	}
	// here we are killing the initial branch
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
	branchNameToKill gitdomain.LocalBranchName
	branches         gitdomain.BranchInfos
	initialBranch    gitdomain.LocalBranchName
	mainBranch       gitdomain.LocalBranchName
	previousBranch   Option[gitdomain.LocalBranchName]
}

func validateKillData(data killData) error {
	switch data.branchToKillType {
	case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
		return nil
	case configdomain.BranchTypeMainBranch:
		return errors.New(messages.KillCannotKillMainBranch)
	case configdomain.BranchTypePerennialBranch:
		return errors.New(messages.KillCannotKillPerennialBranches)
	}
	panic(fmt.Sprintf("unhandled branch type: %s", data.branchToKillType))
}
