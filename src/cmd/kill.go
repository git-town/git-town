package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/sync"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
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

func executeKill(args []string, dryRun, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, initialBranchesSnapshot, initialStashSize, exit, err := determineKillData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateKillData(data)
	if err != nil {
		return err
	}
	steps, finalUndoProgram := killProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               "kill",
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            steps,
		FinalUndoProgram:      finalUndoProgram,
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Config:                  data.config,
		Connector:               nil,
		DialogTestInputs:        &data.dialogTestInputs,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     data.runner,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type killData struct {
	branchNameToKill gitdomain.BranchInfo
	branchTypeToKill configdomain.BranchType
	branchWhenDone   gitdomain.LocalBranchName
	config           configdomain.ValidatedConfig
	dialogTestInputs components.TestInputs
	dryRun           bool
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	parentBranch     Option[gitdomain.LocalBranchName]
	previousBranch   gitdomain.LocalBranchName
	runner           *git.ProdRunner
}

func determineKillData(args []string, repo *execute.OpenRepoResult, dryRun, verbose bool) (*killData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	runner := git.ProdRunner{
		Backend:         repo.Backend,
		CommandsCounter: repo.CommandsCounter,
		Config:          repo.Config,
		FinalMessages:   repo.FinalMessages,
		Frontend:        repo.Frontend,
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                &repo.UnvalidatedConfig.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		HandleUnfinishedState: false,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	branchNameToKill := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchesToKill := gitdomain.LocalBranchNames{branchNameToKill}
	branchToKill, hasBranchToKill := branchesSnapshot.Branches.FindByLocalName(branchNameToKill).Get()
	if !hasBranchToKill {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToKill)
	}
	if branchToKill.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return nil, branchesSnapshot, stashSize, exit, fmt.Errorf(messages.KillBranchOtherWorktree, branchNameToKill)
	}
	branchTypeToKill := repo.UnvalidatedConfig.Config.BranchType(branchNameToKill)
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	var branchWhenDone gitdomain.LocalBranchName
	validatedConfig, err := validate.Config(repo.UnvalidatedConfig, branchesToKill, branchesSnapshot.Branches.LocalBranches(), &repo.Backend, &dialogTestInputs)
	if err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	if branchNameToKill == branchesSnapshot.Active {
		if previousBranch == branchesSnapshot.Active {
			branchWhenDone = validatedConfig.Config.MainBranch
		} else {
			branchWhenDone = previousBranch
		}
	} else {
		branchWhenDone = branchesSnapshot.Active
	}
	parentBranch := validatedConfig.Config.Lineage.Parent(branchToKill.LocalName)
	runner := git.ProdRunner{
		Config:          validatedConfig,
		Backend:         repo.Backend,
		Frontend:        repo.Frontend,
		CommandsCounter: repo.CommandsCounter,
		FinalMessages:   &repo.FinalMessages,
	}
	return &killData{
		config:           validatedConfig.Config,
		branchNameToKill: branchToKill,
		branchTypeToKill: branchTypeToKill,
		branchWhenDone:   branchWhenDone,
		dialogTestInputs: dialogTestInputs,
		dryRun:           dryRun,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    branchesSnapshot.Active,
		parentBranch:     parentBranch,
		previousBranch:   previousBranch,
		runner:           &runner,
	}, branchesSnapshot, stashSize, false, nil
}

func killProgram(config *killData) (runProgram, finalUndoProgram program.Program) {
	prog := program.Program{}
	switch config.branchTypeToKill {
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch:
		killFeatureBranch(&prog, &finalUndoProgram, config)
	case configdomain.BranchTypeObservedBranch, configdomain.BranchTypeContributionBranch:
		killLocalBranch(&prog, &finalUndoProgram, config)
	case configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		panic(fmt.Sprintf("this branch type should have been filtered in validation: %s", config.branchTypeToKill))
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.initialBranch != config.branchNameToKill.LocalName && config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch, config.initialBranch},
	})
	return prog, finalUndoProgram
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killFeatureBranch(prog *program.Program, finalUndoProgram *program.Program, config *killData) {
	if config.branchNameToKill.HasTrackingBranch() && config.config.IsOnline() {
		prog.Add(&opcodes.DeleteTrackingBranch{Branch: config.branchNameToKill.RemoteName})
	}
	killLocalBranch(prog, finalUndoProgram, config)
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killLocalBranch(prog *program.Program, finalUndoProgram *program.Program, config *killData) {
	if config.initialBranch == config.branchNameToKill.LocalName {
		if config.hasOpenChanges {
			prog.Add(&opcodes.CommitOpenChanges{})
			// update the registered initial SHA for this branch so that undo restores the just committed changes
			prog.Add(&opcodes.UpdateInitialBranchLocalSHA{Branch: config.initialBranch})
			// when undoing, manually undo the just committed changes so that they are uncommitted again
			finalUndoProgram.Add(&opcodes.Checkout{Branch: config.branchNameToKill.LocalName})
			finalUndoProgram.Add(&opcodes.UndoLastCommit{})
		}
		prog.Add(&opcodes.Checkout{Branch: config.branchWhenDone})
	}
	prog.Add(&opcodes.DeleteLocalBranch{Branch: config.branchNameToKill.LocalName})
	if parentBranch, hasParentBranch := config.parentBranch.Get(); hasParentBranch && !config.dryRun {
		sync.RemoveBranchFromLineage(sync.RemoveBranchFromLineageArgs{
			Branch:  config.branchNameToKill.LocalName,
			Lineage: config.config.Lineage,
			Parent:  parentBranch,
			Program: prog,
		})
	}
}

func validateKillData(data *killData) error {
	switch data.branchTypeToKill {
	case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch:
		return nil
	case configdomain.BranchTypeMainBranch:
		return errors.New(messages.KillCannotKillMainBranch)
	case configdomain.BranchTypePerennialBranch:
		return errors.New(messages.KillCannotKillPerennialBranches)
	}
	panic(fmt.Sprintf("unhandled branch type: %s", data.branchTypeToKill))
}
