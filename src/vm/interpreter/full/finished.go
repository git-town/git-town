package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v13/src/cli/print"
	"github.com/git-town/git-town/v13/src/config/gitconfig"
	"github.com/git-town/git-town/v13/src/messages"
	"github.com/git-town/git-town/v13/src/undo/undoconfig"
	"github.com/git-town/git-town/v13/src/vm/statefile"
)

// finished is called when executing all steps has successfully finished.
func finished(args ExecuteArgs) error {
	if args.RunState.IsUndo {
		return finishedUndoCommand(args)
	}
	var err error
	args.RunState.EndBranchesSnapshot, err = args.Run.Backend.BranchesSnapshot()
	if err != nil {
		return err
	}
	configGitAccess := gitconfig.Access{Runner: args.Run.Backend.Runner}
	globalSnapshot, _, err := configGitAccess.LoadGlobal(false)
	if err != nil {
		return err
	}
	localSnapshot, _, err := configGitAccess.LoadLocal(false)
	if err != nil {
		return err
	}
	args.RunState.EndConfigSnapshot = undoconfig.ConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	}
	args.RunState.EndStashSize, err = args.Run.Backend.StashSize()
	if err != nil {
		return err
	}
	args.RunState.MarkAsFinished()
	if args.RunState.DryRun {
		return finishedDryRunCommand(args)
	}
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
