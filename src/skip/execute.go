package skip

import (
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v12/src/undo/undobranches"
	fullInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/full"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

func Execute(args ExecuteArgs) error {
	// abort the current op
	for _, opcode := range args.RunState.AbortProgram {
		err := opcode.Run(shared.RunArgs{
			Connector:                       nil,
			DialogTestInputs:                nil,
			Lineage:                         args.Runner.Lineage,
			PrependOpcodes:                  nil,
			RegisterUndoablePerennialCommit: nil,
			Runner:                          args.Runner,
			UpdateInitialBranchLocalSHA:     nil,
		})
		if err != nil {
			panic(err.Error())
		}
	}
	// undo the changes to the current branch
	spans := undobranches.BranchSpans{
		undobranches.BranchSpan{
			Before: *args.RunState.BeforeBranchesSnapshot.Branches.FindByLocalName(args.CurrentBranch),
			After:  *args.RunState.AfterBranchesSnapshot.Branches.FindByLocalName(args.CurrentBranch),
		},
	}
	changes := spans.Changes()
	undoCurrentBranchProgram := changes.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		Config:                   &args.Runner.FullConfig,
		FinalBranch:              args.CurrentBranch,
		InitialBranch:            args.CurrentBranch,
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	for _, opcode := range undoCurrentBranchProgram {
		err := opcode.Run(shared.RunArgs{
			Connector:                       nil,
			DialogTestInputs:                nil,
			Lineage:                         args.Runner.Lineage,
			PrependOpcodes:                  nil,
			RegisterUndoablePerennialCommit: nil,
			Runner:                          args.Runner,
			UpdateInitialBranchLocalSHA:     nil,
		})
		if err != nil {
			panic(err.Error())
		}
	}

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
		InitialBranchesSnapshot: args.RunState.BeforeBranchesSnapshot,
		InitialConfigSnapshot:   args.RunState.BeforeConfigSnapshot,
		InitialStashSize:        args.RunState.BeforeStashSize,
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
