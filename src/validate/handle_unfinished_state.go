package validate

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/dialog"
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	fullInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/full"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/git-town/git-town/v12/src/vm/statefile"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
func HandleUnfinishedState(args UnfinishedStateArgs) (quit bool, err error) {
	runState, err := statefile.Load(args.RootDir)
	if err != nil {
		return false, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || runState.IsFinished() {
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
		return abortRunstate(runState, args)
	case dialog.ResponseSkip:
		return skipRunstate(runState, args)
	case dialog.ResponseQuit:
		return true, nil
	}
	return false, fmt.Errorf(messages.DialogUnexpectedResponse, response)
}

type UnfinishedStateArgs struct {
	Connector               hostingdomain.Connector
	DialogTestInputs        components.TestInputs
	InitialBranchesSnapshot gitdomain.BranchesSnapshot
	InitialConfigSnapshot   undoconfig.ConfigSnapshot
	InitialStashSize        gitdomain.StashSize
	Lineage                 configdomain.Lineage
	PushHook                configdomain.PushHook
	RootDir                 gitdomain.RepoRootDir
	Run                     *git.ProdRunner
	Verbose                 bool
}

func abortRunstate(runState *runstate.RunState, args UnfinishedStateArgs) (bool, error) {
	abortRunState := runState.CreateAbortRunState()
	return true, fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               args.Connector,
		DialogTestInputs:        &args.DialogTestInputs,
		FullConfig:              &args.Run.FullConfig,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSize:        args.InitialStashSize,
		RootDir:                 args.RootDir,
		Run:                     args.Run,
		RunState:                &abortRunState,
		Verbose:                 args.Verbose,
	})
}

func continueRunstate(runState *runstate.RunState, args UnfinishedStateArgs) (bool, error) {
	repoStatus, err := args.Run.Backend.RepoStatus()
	if err != nil {
		return false, err
	}
	if repoStatus.Conflicts {
		return false, errors.New(messages.ContinueUnresolvedConflicts)
	}
	return true, fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               args.Connector,
		DialogTestInputs:        &args.DialogTestInputs,
		FullConfig:              &args.Run.FullConfig,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSize:        args.InitialStashSize,
		RootDir:                 args.RootDir,
		Run:                     args.Run,
		RunState:                runState,
		Verbose:                 args.Verbose,
	})
}

func discardRunstate(rootDir gitdomain.RepoRootDir) (bool, error) {
	err := statefile.Delete(rootDir)
	return false, err
}

func skipRunstate(runState *runstate.RunState, args UnfinishedStateArgs) (bool, error) {
	skipRunState := runState.CreateSkipRunState()
	return true, fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               args.Connector,
		DialogTestInputs:        &args.DialogTestInputs,
		FullConfig:              &args.Run.FullConfig,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSize:        args.InitialStashSize,
		RootDir:                 args.RootDir,
		Run:                     args.Run,
		RunState:                &skipRunState,
		Verbose:                 args.Verbose,
	})
}
