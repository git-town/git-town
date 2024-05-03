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
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/skip"
	"github.com/git-town/git-town/v14/src/undo"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/git-town/git-town/v14/src/vm/statefile"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
func HandleUnfinishedState(args UnfinishedStateArgs) (quit bool, err error) {
	runStateOpt, err := statefile.Load(args.RootDir)
	if err != nil {
		return false, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	runState, hasRunState := runStateOpt.Get()
	if !hasRunState || runState.IsFinished() {
		return false, nil
	}
	response, aborted, err := dialog.AskHowToHandleUnfinishedRunState(
		runState.Command,
		runState.UnfinishedDetails.EndBranch,
		runState.UnfinishedDetails.EndTime,
		runState.UnfinishedDetails.CanSkip,
		args.DialogTestInputs.Next(),
	)
	if err != nil {
		return quit, err
	}
	if aborted {
		return quit, errors.New("user aborted")
	}
	switch response {
	case dialog.ResponseDiscard:
		return discardRunstate(args.RootDir)
	case dialog.ResponseContinue:
		return continueRunstate(runState, args)
	case dialog.ResponseUndo:
		return true, undo.Execute(undo.ExecuteArgs{
			Config:           args.Config.Config,
			HasOpenChanges:   args.HasOpenChanges,
			InitialStashSize: args.InitialStashSize,
			Lineage:          args.Lineage,
			RootDir:          args.RootDir,
			RunState:         runState,
			Verbose:          args.Verbose,
		})
	case dialog.ResponseSkip:
		return true, skip.Execute(skip.ExecuteArgs{
			Connector:      args.Connector,
			CurrentBranch:  args.CurrentBranch,
			HasOpenChanges: args.HasOpenChanges,
			RootDir:        args.RootDir,
			RunState:       runState,
			TestInputs:     args.DialogTestInputs,
			Verbose:        args.Verbose,
		})
	case dialog.ResponseQuit:
		return true, nil
	}
	return false, fmt.Errorf(messages.DialogUnexpectedResponse, response)
}

type UnfinishedStateArgs struct {
	Backend                 git.BackendCommands
	Config                  config.Config
	Connector               hostingdomain.Connector
	CurrentBranch           gitdomain.LocalBranchName
	DialogTestInputs        components.TestInputs
	HasOpenChanges          bool
	InitialBranchesSnapshot gitdomain.BranchesSnapshot
	InitialConfigSnapshot   undoconfig.ConfigSnapshot
	InitialStashSize        gitdomain.StashSize
	Lineage                 configdomain.Lineage
	PushHook                configdomain.PushHook
	RootDir                 gitdomain.RepoRootDir
	Verbose                 bool
}

func continueRunstate(runState runstate.RunState, args UnfinishedStateArgs) (bool, error) {
	repoStatus, err := args.Backend.RepoStatus()
	if err != nil {
		return false, err
	}
	if repoStatus.Conflicts {
		return false, errors.New(messages.ContinueUnresolvedConflicts)
	}
	return true, fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Config:                  args.Config,
		Connector:               args.Connector,
		DialogTestInputs:        &args.DialogTestInputs,
		HasOpenChanges:          repoStatus.OpenChanges,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSize:        args.InitialStashSize,
		RootDir:                 args.RootDir,
		RunState:                runState,
		Verbose:                 args.Verbose,
	})
}

func discardRunstate(rootDir gitdomain.RepoRootDir) (bool, error) {
	err := statefile.Delete(rootDir)
	return false, err
}
