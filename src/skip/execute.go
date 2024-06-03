package skip

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undobranches"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	lightInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/light"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// executes the "skip" command at the given runstate
func Execute(args ExecuteArgs) error {
	lightInterpreter.Execute(lightInterpreter.ExecuteArgs{
		Backend:       args.Backend,
		Config:        args.Config,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Prog:          args.RunState.AbortProgram,
	})
	err := revertChangesToCurrentBranch(args)
	if err != nil {
		return err
	}
	args.RunState.RunProgram = removeOpcodesForCurrentBranch(args.RunState.RunProgram)
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 args.Backend,
		CommandsCounter:         args.CommandsCounter,
		Config:                  args.Config,
		Connector:               args.Connector,
		DialogTestInputs:        args.TestInputs,
		FinalMessages:           args.FinalMessages,
		Frontend:                args.Frontend,
		HasOpenChanges:          args.HasOpenChanges,
		InitialBranch:           args.InitialBranch,
		InitialBranchesSnapshot: args.RunState.BeginBranchesSnapshot,
		InitialConfigSnapshot:   args.RunState.BeginConfigSnapshot,
		InitialStashSize:        args.RunState.BeginStashSize,
		RootDir:                 args.RootDir,
		RunState:                args.RunState,
		Verbose:                 args.Verbose,
	})
}

type ExecuteArgs struct {
	Backend         git.BackendCommands
	CommandsCounter gohacks.Counter
	Config          config.ValidatedConfig
	Connector       Option[hostingdomain.Connector]
	FinalMessages   stringslice.Collector
	Frontend        git.FrontendCommands
	HasOpenChanges  bool
	InitialBranch   gitdomain.LocalBranchName
	RootDir         gitdomain.RepoRootDir
	RunState        runstate.RunState
	TestInputs      components.TestInputs
	Verbose         bool
}

// removes the remaining opcodes for the current branch from the given program
func removeOpcodesForCurrentBranch(prog program.Program) program.Program {
	result := make(program.Program, 0, len(prog)-1)
	skipping := true
	for _, opcode := range prog {
		if shared.IsEndOfBranchProgramOpcode(opcode) {
			skipping = false
			continue
		}
		if !skipping {
			result.Add(opcode)
		}
	}
	return result
}

func revertChangesToCurrentBranch(args ExecuteArgs) error {
	before := args.RunState.BeginBranchesSnapshot.Branches.FindByLocalName(args.InitialBranch)
	if before.IsNone() {
		return fmt.Errorf(messages.SkipNoInitialBranchInfo, args.InitialBranch)
	}
	afterSnapshot, hasAfterSnapshot := args.RunState.EndBranchesSnapshot.Get()
	if !hasAfterSnapshot {
		return errors.New(messages.SkipNoFinalSnapshot)
	}
	spans := undobranches.BranchSpans{
		undobranches.BranchSpan{
			Before: before,
			After:  afterSnapshot.Branches.FindByLocalName(args.InitialBranch),
		},
	}
	undoCurrentBranchProgram := spans.Changes().UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		BeginBranch:              args.InitialBranch,
		Config:                   args.Config.Config,
		EndBranch:                args.InitialBranch,
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	lightInterpreter.Execute(lightInterpreter.ExecuteArgs{
		Backend:       args.Backend,
		Config:        args.Config,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Prog:          undoCurrentBranchProgram,
	})
	return nil
}
