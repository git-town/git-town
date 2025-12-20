package skip

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/internal/undo/undobranches"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/lightinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ExecuteArgs struct {
	Backend         subshelldomain.RunnerQuerier
	CommandsCounter Mutable[gohacks.Counter]
	Config          config.ValidatedConfig
	Connector       Option[forgedomain.Connector]
	FinalMessages   stringslice.Collector
	Frontend        subshelldomain.Runner
	Git             git.Commands
	HasOpenChanges  bool
	InitialBranch   gitdomain.LocalBranchName
	Inputs          dialogcomponents.Inputs
	Park            configdomain.Park
	RootDir         gitdomain.RepoRootDir
	RunState        runstate.RunState
}

// executes the "skip" command at the given runstate
func Execute(args ExecuteArgs) error {
	skipProgram := args.RunState.AbortProgram
	if args.Park {
		skipProgram = append(skipProgram, &opcodes.BranchTypeOverrideSet{
			Branch:     args.InitialBranch,
			BranchType: configdomain.BranchTypeParkedBranch,
		})
	}
	lightinterpreter.Execute(lightinterpreter.ExecuteArgs{
		Backend:       args.Backend,
		BranchInfos:   args.RunState.BeginBranchesSnapshot.Branches,
		Config:        args.Config,
		Connector:     args.Connector,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Git:           args.Git,
		Prog:          skipProgram,
	})
	args.RunState.AbortProgram = program.Program{}
	if err := revertChangesToCurrentBranch(args); err != nil {
		return err
	}
	args.RunState.RunProgram = RemoveOpcodesForCurrentBranch(args.RunState.RunProgram)
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 args.Backend,
		CommandsCounter:         args.CommandsCounter,
		Config:                  args.Config,
		Connector:               args.Connector,
		FinalMessages:           args.FinalMessages,
		Frontend:                args.Frontend,
		Git:                     args.Git,
		HasOpenChanges:          args.HasOpenChanges,
		InitialBranch:           args.InitialBranch,
		InitialBranchesSnapshot: args.RunState.BeginBranchesSnapshot,
		InitialConfigSnapshot:   args.RunState.BeginConfigSnapshot,
		InitialStashSize:        args.RunState.BeginStashSize,
		Inputs:                  args.Inputs,
		PendingCommand:          Some(args.RunState.Command),
		RootDir:                 args.RootDir,
		RunState:                args.RunState,
	})
}

// removes the remaining opcodes for the current branch from the given program
func RemoveOpcodesForCurrentBranch(prog program.Program) program.Program {
	result := make(program.Program, 0, len(prog)-1)
	skipping := true
	for _, opcode := range prog {
		if opcodes.IsEndOfBranchProgramOpcode(opcode) && skipping {
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
			Before: before.ToOption(),
			After:  afterSnapshot.Branches.FindByLocalName(args.InitialBranch).ToOption(),
		},
	}
	undoCurrentBranchProgram := spans.Changes().UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		BeginBranch:              args.InitialBranch,
		Config:                   args.Config,
		EndBranch:                args.InitialBranch,
		FinalMessages:            args.FinalMessages,
		UndoAPIProgram:           args.RunState.UndoAPIProgram,
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	lightinterpreter.Execute(lightinterpreter.ExecuteArgs{
		Backend:       args.Backend,
		BranchInfos:   args.RunState.BeginBranchesSnapshot.Branches,
		Config:        args.Config,
		Connector:     args.Connector,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Git:           args.Git,
		Prog:          undoCurrentBranchProgram,
	})
	return nil
}
