package skip

import (
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v12/src/undo/undobranches"
	fullInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/full"
	lightInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/light"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// executes the "skip" command at the given runstate
func Execute(args ExecuteArgs) error {
	// abort the current op
	lightInterpreter.Execute(args.RunState.AbortProgram, args.Runner, args.Runner.Lineage)
	// undo the changes to the current branch
	spans := undobranches.BranchSpans{
		undobranches.BranchSpan{
			Before: *args.RunState.BeginBranchesSnapshot.Branches.FindByLocalName(args.CurrentBranch),
			After:  *args.RunState.EndBranchesSnapshot.Branches.FindByLocalName(args.CurrentBranch),
		},
	}
	changes := spans.Changes()
	undoCurrentBranchProgram := changes.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		BeginBranch:              args.CurrentBranch,
		Config:                   &args.Runner.FullConfig,
		EndBranch:                args.CurrentBranch,
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	lightInterpreter.Execute(undoCurrentBranchProgram, args.Runner, args.Runner.Lineage)
	// remove the remaining opcodes for the current branch from the program
	newProgram := program.Program{}
	skipping := true
	for _, opcode := range args.RunState.RunProgram {
		if shared.IsEndOfBranchProgramOpcode(opcode) {
			skipping = false
		}
		if !skipping {
			newProgram.Add(opcode)
		}
	}
	newProgram.MoveToEnd(&opcodes.RestoreOpenChanges{})
	args.RunState.RunProgram = newProgram
	// continue executing the program
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               args.Connector,
		DialogTestInputs:        &args.TestInputs,
		FullConfig:              &args.Runner.FullConfig,
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
