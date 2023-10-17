package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli/dialog"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/git-town/git-town/v9/src/vm/interpreter"
	"github.com/git-town/git-town/v9/src/vm/persistence"
	"github.com/git-town/git-town/v9/src/vm/runstate"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
func HandleUnfinishedState(args UnfinishedStateArgs) (quit bool, err error) {
	runState, err := persistence.Load(args.RootDir)
	if err != nil {
		return false, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return false, nil
	}
	response, err := dialog.AskHowToHandleUnfinishedRunState(
		runState.Command,
		runState.UnfinishedDetails.EndBranch,
		runState.UnfinishedDetails.EndTime,
		runState.UnfinishedDetails.CanSkip,
	)
	if err != nil {
		return quit, err
	}
	switch response {
	case dialog.ResponseDiscard:
		return discardRunstate(args.RootDir)
	case dialog.ResponseContinue:
		return continueRunstate(runState, args)
	case dialog.ResponseAbort:
		return abortRunstate(runState, args)
	case dialog.ResponseSkip:
		return skipRunstate(runState, args)
	case dialog.ResponseQuit:
		return true, nil
	default:
		return false, fmt.Errorf(messages.DialogUnexpectedResponse, response)
	}
}

type UnfinishedStateArgs struct {
	Connector               hosting.Connector
	Debug                   bool
	Lineage                 config.Lineage
	InitialBranchesSnapshot domain.BranchesSnapshot
	InitialConfigSnapshot   undo.ConfigSnapshot
	InitialStashSnapshot    domain.StashSnapshot
	PushHook                bool
	RootDir                 domain.RepoRootDir
	Run                     *git.ProdRunner
}

func abortRunstate(runState *runstate.RunState, args UnfinishedStateArgs) (bool, error) {
	abortRunState := runState.CreateAbortRunState()
	return true, interpreter.Execute(interpreter.ExecuteArgs{
		Connector:               args.Connector,
		Debug:                   args.Debug,
		Lineage:                 args.Lineage,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSnapshot:    args.InitialStashSnapshot,
		NoPushHook:              !args.PushHook,
		RootDir:                 args.RootDir,
		Run:                     args.Run,
		RunState:                &abortRunState,
	})
}

func continueRunstate(runState *runstate.RunState, args UnfinishedStateArgs) (bool, error) {
	repoStatus, err := args.Run.Backend.RepoStatus()
	if err != nil {
		return false, err
	}
	if repoStatus.Conflicts {
		return false, fmt.Errorf(messages.ContinueUnresolvedConflicts)
	}
	return true, interpreter.Execute(interpreter.ExecuteArgs{
		Connector:               args.Connector,
		Debug:                   args.Debug,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSnapshot:    args.InitialStashSnapshot,
		Lineage:                 args.Lineage,
		NoPushHook:              !args.PushHook,
		RootDir:                 args.RootDir,
		Run:                     args.Run,
		RunState:                runState,
	})
}

func discardRunstate(rootDir domain.RepoRootDir) (bool, error) {
	err := persistence.Delete(rootDir)
	return false, err
}

func skipRunstate(runState *runstate.RunState, args UnfinishedStateArgs) (bool, error) {
	skipRunState := runState.CreateSkipRunState()
	return true, interpreter.Execute(interpreter.ExecuteArgs{
		Connector:               args.Connector,
		Debug:                   args.Debug,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSnapshot:    args.InitialStashSnapshot,
		Lineage:                 args.Lineage,
		NoPushHook:              !args.PushHook,
		RootDir:                 args.RootDir,
		Run:                     args.Run,
		RunState:                &skipRunState,
	})
}
