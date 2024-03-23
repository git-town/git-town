package skip

import (
	"github.com/git-town/git-town/v13/src/cli/dialog/components"
	"github.com/git-town/git-town/v13/src/git"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v13/src/undo/undobranches"
	fullInterpreter "github.com/git-town/git-town/v13/src/vm/interpreter/full"
	lightInterpreter "github.com/git-town/git-town/v13/src/vm/interpreter/light"
	"github.com/git-town/git-town/v13/src/vm/program"
	"github.com/git-town/git-town/v13/src/vm/runstate"
	"github.com/git-town/git-town/v13/src/vm/shared"
)

// executes the "skip" command at the given runstate
func Execute(args ExecuteArgs) error {
	lightInterpreter.Execute(args.RunState.AbortProgram, args.Runner, args.Runner.Config.FullConfig.Lineage)
	revertChangesToCurrentBranch(args)
	args.RunState.RunProgram = removeOpcodesForCurrentBranch(args.RunState.RunProgram)
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               args.Connector,
		DialogTestInputs:        &args.TestInputs,
		FullConfig:              &args.Runner.Config.FullConfig,
		HasOpenChanges:          args.HasOpenChanges,
		InitialBranchesSnapshot: args.RunState.BeginBranchesSnapshot,
		InitialConfigSnapshot:   args.RunState.BeginConfigSnapshot,
		InitialStashSize:        args.RunState.BeginStashSize,
		RootDir:                 args.RootDir,
		Run:                     args.Runner,
		RunState:                args.RunState,
		Verbose:                 args.Verbose,
	})
}

type ExecuteArgs struct {
	Connector      hostingdomain.Connector
	CurrentBranch  gitdomain.LocalBranchName
	HasOpenChanges bool
	RootDir        gitdomain.RepoRootDir
	RunState       *runstate.RunState
	Runner         *git.ProdRunner
	TestInputs     components.TestInputs
	Verbose        bool
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

func revertChangesToCurrentBranch(args ExecuteArgs) {
	spans := undobranches.BranchSpans{
		undobranches.BranchSpan{
			Before: *args.RunState.BeginBranchesSnapshot.Branches.FindByLocalName(args.CurrentBranch),
			After:  *args.RunState.EndBranchesSnapshot.Branches.FindByLocalName(args.CurrentBranch),
		},
	}
	undoCurrentBranchProgram := spans.Changes().UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		BeginBranch:              args.CurrentBranch,
		Config:                   &args.Runner.Config.FullConfig,
		EndBranch:                args.CurrentBranch,
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	lightInterpreter.Execute(undoCurrentBranchProgram, args.Runner, args.Runner.Config.FullConfig.Lineage)
}
