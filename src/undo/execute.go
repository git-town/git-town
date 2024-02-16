package undo

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo/undobranches"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	"github.com/git-town/git-town/v12/src/undo/undostash"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/git-town/git-town/v12/src/vm/shared"
	"github.com/git-town/git-town/v12/src/vm/statefile"
)

// undoes the persisted runstate
func Execute(args ExecuteArgs) error {
	undoProgram := program.Program{}
	undoProgram.AddProgram(args.RunState.AbortProgram)

	// undo branch changes
	branchSpans := undobranches.NewBranchSpans(args.RunState.BeforeBranchesSnapshot, args.RunState.AfterBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	undoBranchesProgram := branchChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		Config:                   args.FullConfig,
		FinalBranch:              args.RunState.AfterBranchesSnapshot.Active,
		InitialBranch:            args.RunState.BeforeBranchesSnapshot.Active,
		UndoablePerennialCommits: []gitdomain.SHA{},
	})
	undoProgram.AddProgram(undoBranchesProgram)

	// undo config changes
	fmt.Println("111111111", args.RunState.BeforeConfigSnapshot)
	fmt.Println("222222222", args.RunState.AfterConfigSnapshot)
	configSpans := undoconfig.NewConfigDiffs(args.RunState.BeforeConfigSnapshot, args.RunState.AfterConfigSnapshot)
	configUndoProgram := configSpans.UndoProgram()
	undoProgram.AddProgram(configUndoProgram)

	// undo stash changes
	stashDiff := undostash.NewStashDiff(args.RunState.BeforeStashSize, args.InitialStashSize)
	undoStashProgram := stashDiff.Program()
	undoProgram.AddProgram(undoStashProgram)

	undoProgram.AddProgram(args.RunState.FinalUndoProgram)
	undoProgram.Add(&opcodes.Checkout{Branch: args.RunState.BeforeBranchesSnapshot.Active})
	undoProgram.RemoveDuplicateCheckout()

	cmdhelpers.Wrap(&undoProgram, cmdhelpers.WrapOptions{
		DryRun:                   args.RunState.DryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         args.HasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{args.RunState.BeforeBranchesSnapshot.Active},
	})

	// execute the undo program
	for _, opcode := range undoProgram {
		err := opcode.Run(shared.RunArgs{
			Connector:                       nil,
			DialogTestInputs:                nil,
			Lineage:                         args.Lineage,
			PrependOpcodes:                  nil,
			RegisterUndoablePerennialCommit: nil,
			Runner:                          args.Runner,
			UpdateInitialBranchLocalSHA:     nil,
		})
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
	}

	err := statefile.Delete(args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateDeleteProblem, err)
	}
	print.Footer(args.Verbose, args.Runner.CommandsCounter.Count(), args.Runner.FinalMessages.Result())
	return nil
}

type ExecuteArgs struct {
	FullConfig       *configdomain.FullConfig
	HasOpenChanges   bool
	InitialStashSize gitdomain.StashSize
	Lineage          configdomain.Lineage
	RootDir          gitdomain.RepoRootDir
	RunState         runstate.RunState
	Runner           *git.ProdRunner
	Verbose          bool
}
