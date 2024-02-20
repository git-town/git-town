package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/config/gitconfig"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	"github.com/git-town/git-town/v12/src/vm/statefile"
)

// finished is called when executing all steps has successfully finished.
func finished(args ExecuteArgs) error {
	if args.RunState.IsUndo {
		return finishedUndoCommand(args)
	}
	var err error
	args.RunState.AfterBranchesSnapshot, err = args.Run.Backend.BranchesSnapshot()
	if err != nil {
		return err
	}
	configGitAccess := gitconfig.Access{Runner: args.Run.Backend}
	globalSnapshot, _, err := configGitAccess.LoadGlobal()
	if err != nil {
		return err
	}
	localSnapshot, _, err := configGitAccess.LoadLocal()
	if err != nil {
		return err
	}
	args.RunState.AfterConfigSnapshot = undoconfig.ConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	}
	args.RunState.MarkAsFinished()
	if args.RunState.DryRun {
		return finishedDryRunCommand(args)
	}
	undoProgram, err := undo.CreateUndoProgram(undo.CreateUndoProgramArgs{
		DryRun:                   args.RunState.DryRun,
		FinalBranchesSnapshot:    args.RunState.AfterBranchesSnapshot,
		FinalConfigSnapshot:      args.RunState.AfterConfigSnapshot,
		InitialBranchesSnapshot:  args.InitialBranchesSnapshot,
		InitialConfigSnapshot:    args.InitialConfigSnapshot,
		InitialStashSize:         args.InitialStashSize,
		NoPushHook:               args.NoPushHook(),
		Run:                      args.Run,
		RunState:                 *args.RunState,
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	if err != nil {
		return err
	}
	undoProgram.AddProgram(args.RunState.FinalUndoProgram)
	args.RunState.UndoProgram = undoProgram
	err = statefile.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Verbose, args.Run.CommandsCounter.Count(), args.Run.FinalMessages.Result())
	return nil
}

func finishedDryRunCommand(args ExecuteArgs) error {
	args.RunState.MarkAsFinished()
	err := statefile.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Verbose, args.Run.CommandsCounter.Count(), args.Run.FinalMessages.Result())
	return nil
}

func finishedUndoCommand(args ExecuteArgs) error {
	err := statefile.Delete(args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateDeleteProblem, err)
	}
	print.Footer(args.Verbose, args.Run.CommandsCounter.Count(), args.Run.FinalMessages.Result())
	return nil
}
