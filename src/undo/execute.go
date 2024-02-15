package undo

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo/undobranches"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	"github.com/git-town/git-town/v12/src/undo/undostash"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/shared"
	"github.com/git-town/git-town/v12/src/vm/statefile"
)

func Execute(rootDir gitdomain.RepoRootDir) error {
	runState, err := statefile.Load(rootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	undoProgram := program.Program{}
	undoProgram.AddProgram(runState.AbortProgram)

	// undo branch changes
	branchSpans := undobranches.NewBranchSpans(runState.BeforeBranchesSnapshot, runState.AfterBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	branchUndoProgram := branchChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		Config:                   config.FullConfig,
		FinalBranch:              runState.AfterBranchesSnapshot.Active,
		InitialBranch:            runState.BeforeBranchesSnapshot.Active,
		UndoablePerennialCommits: []gitdomain.SHA{},
	})
	undoProgram.AddProgram(branchUndoProgram)

	// undo config changes
	configSpans := undoconfig.NewConfigDiffs(runState.BeforeConfigSnapshot, runState.AfterConfigSnapshot)
	configUndoProgram := configSpans.UndoProgram()
	undoProgram.AddProgram(configUndoProgram)

	// undo stash changes
	stashDiff := undostash.NewStashDiff(runState.BeforeStashSize, initialStashSize)
	undoStashProgram := stashDiff.Program()
	undoProgram.AddProgram(undoStashProgram)

	undoProgram.AddProgram(runState.FinalUndoProgram)

	cmdhelpers.Wrap(&undoProgram, cmdhelpers.WrapOptions{
		DryRun:                   runState.DryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch},
	})

	// execute the undo program
	for _, opcode := range branchUndoProgram {
		err := opcode.Run(shared.RunArgs{
			Connector:                       nil,
			DialogTestInputs:                nil,
			Lineage:                         config.Lineage,
			PrependOpcodes:                  nil,
			RegisterUndoablePerennialCommit: nil,
			Runner:                          repo.Runner,
			UpdateInitialBranchLocalSHA:     nil,
		})
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
	}
	return nil
}
