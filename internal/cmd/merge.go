package cmd

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	fullInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const mergeDesc = "Merges the current branch in a stack with its parent"

const mergeHelp = `
Merges the current branch with its parent branch.
Both branches must be feature branches and be in sync with their tracking branches.
`

func mergeCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     "merge",
		Args:    cobra.NoArgs,
		GroupID: "refactoring",
		Short:   mergeDesc,
		Long:    cmdhelpers.Long(mergeDesc, mergeHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeMerge(readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeMerge(dryRun configdomain.DryRun, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, err := determineMergeData(repo)
	if err != nil {
		return err
	}
	err = validateMergeData(data)
	if err != nil {
		return err
	}
	runProgram := mergeProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "prepend",
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
	branchesSnapshot gitdomain.BranchesSnapshot
	config           config.ValidatedConfig
	connector        Option[hostingdomain.Connector]
	dialogTestInputs components.TestInputs
	dryRun           configdomain.DryRun
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	parentBranch     gitdomain.BranchName
	previousBranch   Option[gitdomain.LocalBranchName]
	stashSize        gitdomain.StashSize
}

func determineMergeData(repo execute.OpenRepoResult) (mergeData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return mergeData{}, err
	}
	return mergeData{
		branchesSnapshot: branchesSnapshot,
	}, err
}

func mergeProgram(data mergeData) program.Program {
	prog := NewMutable(&program.Program{})
	switch data.config.NormalConfig.SyncFeatureStrategy {
	case configdomain.SyncFeatureStrategyCompress:
		mergeUsingCompressStrategy(prog, data)
	case configdomain.SyncFeatureStrategyMerge:
		mergeUsingMergeStrategy(prog, data)
	case configdomain.SyncFeatureStrategyRebase:
		mergeUsingRebaseStrategy(prog, data)
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return prog.Get()
}

func mergeUsingCompressStrategy(prog Mutable[program.Program], data mergeData) {}

func mergeUsingMergeStrategy(prog Mutable[program.Program], data mergeData) {
	prog.Value.Add(&opcodes.Merge{
		Branch: data.parentBranch,
	})
}

func mergeUsingRebaseStrategy(prog Mutable[program.Program], data mergeData) {}

func validateMergeData(data mergeData) error {
	if data.hasOpenChanges {
		return errors.New(messages.MergeOpenChanges)
	}
	// ensure branches in sync with tracking branches
	return nil
}
