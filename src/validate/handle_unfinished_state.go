package validate

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/skip"
	"github.com/git-town/git-town/v14/src/undo"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/git-town/git-town/v14/src/vm/statefile"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
func HandleUnfinishedState(args UnfinishedStateArgs) (bool, error) {
	runStateOpt, err := statefile.Load(args.RootDir)
	if err != nil {
		return false, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	runState, hasRunState := runStateOpt.Get()
	if !hasRunState || runState.IsFinished() {
		return false, nil
	}
	response, exit, err := dialog.AskHowToHandleUnfinishedRunState(
		runState.Command,
		runState.UnfinishedDetails.EndBranch,
		runState.UnfinishedDetails.EndTime,
		runState.UnfinishedDetails.CanSkip,
		args.DialogTestInputs.Next(),
	)
	if err != nil {
		return false, err
	}
	if exit {
		return exit, errors.New("user aborted")
	}
	switch response {
	case dialog.ResponseDiscard:
		return discardRunstate(args.RootDir)
	case dialog.ResponseContinue:
		return continueRunstate(runState, args)
	case dialog.ResponseUndo:
		return true, undo.Execute(undo.ExecuteArgs{
			Backend:          args.Backend,
			CommandsCounter:  args.CommandsCounter,
			Config:           args.Config,
			FinalMessages:    args.FinalMessages,
			Frontend:         args.Frontend,
			HasOpenChanges:   args.HasOpenChanges,
			InitialStashSize: runState.BeginStashSize,
			Lineage:          args.Lineage,
			RootDir:          args.RootDir,
			RunState:         runState,
			Verbose:          args.Verbose,
		})
	case dialog.ResponseSkip:
		currentBranch, err := args.Backend.CurrentBranch()
		if err != nil {
			return false, err
		}
		return true, skip.Execute(skip.ExecuteArgs{
			Backend:         args.Backend,
			CommandsCounter: args.CommandsCounter,
			Config:          args.Config,
			Connector:       args.Connector,
			CurrentBranch:   currentBranch,
			FinalMessages:   args.FinalMessages,
			Frontend:        args.Frontend,
			HasOpenChanges:  args.HasOpenChanges,
			RootDir:         args.RootDir,
			RunState:        runState,
			TestInputs:      args.DialogTestInputs,
			Verbose:         args.Verbose,
		})
	case dialog.ResponseQuit:
		return true, nil
	}
	return false, fmt.Errorf(messages.DialogUnexpectedResponse, response)
}

type UnfinishedStateArgs struct {
	Backend          git.BackendCommands
	CommandsCounter  gohacks.Counter
	Config           config.Config
	Connector        hostingdomain.Connector
	DialogTestInputs components.TestInputs
	FinalMessages    stringslice.Collector
	Frontend         git.FrontendCommands
	HasOpenChanges   bool
	Lineage          configdomain.Lineage
	PushHook         configdomain.PushHook
	RepoStatus       gitdomain.RepoStatus
	RootDir          gitdomain.RepoRootDir
	Verbose          bool
}

func continueRunstate(runState runstate.RunState, args UnfinishedStateArgs) (bool, error) {
	if args.RepoStatus.Conflicts {
		return false, errors.New(messages.ContinueUnresolvedConflicts)
	}
	return true, fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 args.Backend,
		CommandsCounter:         args.CommandsCounter,
		Config:                  args.Config,
		Connector:               args.Connector,
		DialogTestInputs:        args.DialogTestInputs,
		FinalMessages:           args.FinalMessages,
		Frontend:                args.Frontend,
		HasOpenChanges:          args.RepoStatus.OpenChanges,
		InitialBranchesSnapshot: runState.BeginBranchesSnapshot,
		InitialConfigSnapshot:   runState.BeginConfigSnapshot,
		InitialStashSize:        runState.BeginStashSize,
		RootDir:                 args.RootDir,
		RunState:                runState,
		Verbose:                 args.Verbose,
	})
}

func discardRunstate(rootDir gitdomain.RepoRootDir) (bool, error) {
	err := statefile.Delete(rootDir)
	return false, err
}
