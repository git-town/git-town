package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo"
	"github.com/git-town/git-town/v12/src/vm/interpreter"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/shared"
	"github.com/git-town/git-town/v12/src/vm/statefile"
	"github.com/spf13/cobra"
)

const skipDesc = "Restarts the last run git-town command by skipping the current branch"

func skipCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "skip",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   skipDesc,
		Long:    cmdhelpers.Long(skipDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSkip(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSkip(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	initialBranchesSnapshot, initialStashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	runState, err := statefile.Load(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || runState.IsFinished() {
		return errors.New(messages.SkipNothingToDo)
	}
	if !runState.UnfinishedDetails.CanSkip {
		return errors.New(messages.SkipBranchHasConflicts)
	}

	// abort the current op
	for _, opcode := range runState.AbortProgram {
		err = opcode.Run(shared.RunArgs{
			Connector:                       nil,
			DialogTestInputs:                nil,
			Lineage:                         repo.Runner.Lineage,
			PrependOpcodes:                  nil,
			RegisterUndoablePerennialCommit: nil,
			Runner:                          repo.Runner,
			UpdateInitialBranchLocalSHA:     nil,
		})
		if err != nil {
			panic(err.Error())
		}
	}
	// undo the changes to the current branch
	undo.Execute(undo.ExecuteArgs{
		FullConfig:       &repo.Runner.FullConfig,
		HasOpenChanges:   false,
		InitialStashSize: initialStashSize,
		Lineage:          repo.Runner.Lineage,
		RootDir:          repo.RootDir,
		Runner:           repo.Runner,
		RunState:         *runState,
	})

	// remove the remaining opcodes for the current branch from the program
	newProgram := program.Program{}
	skipping := true
	for _, opcode := range runState.RunProgram {
		if shared.IsEndOfBranchProgramOpcode(opcode) {
			skipping = false
		}
		if !skipping {
			newProgram.Add(opcode)
		}
	}
	newProgram.MoveToEnd(&opcodes.RestoreOpenChanges{})
	runState.RunProgram = newProgram

	// continue executing the program
	return interpreter.Execute(interpreter.ExecuteArgs{
		FullConfig:              &repo.Runner.FullConfig,
		RunState:                runState,
		Run:                     repo.Runner,
		Connector:               nil,
		DialogTestInputs:        &dialogTestInputs,
		Verbose:                 verbose,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
	})
}
