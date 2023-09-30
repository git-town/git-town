package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/runvm"
	"github.com/git-town/git-town/v9/src/undo"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
func HandleUnfinishedState(args HandleUnfinishedStateArgs) (quit bool, err error) {
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

type HandleUnfinishedStateArgs struct {
	Run                     *git.ProdRunner
	Connector               hosting.Connector
	RootDir                 domain.RepoRootDir
	Lineage                 config.Lineage
	InitialBranchesSnapshot domain.BranchesSnapshot
	InitialConfigSnapshot   undo.ConfigSnapshot
	InitialStashSnapshot    domain.StashSnapshot
	PushHook                bool
}

func abortRunstate(runState *runstate.RunState, args HandleUnfinishedStateArgs) (bool, error) {
	abortRunState := runState.CreateAbortRunState()
	return true, runvm.Execute(runvm.ExecuteArgs{
		RunState:                &abortRunState,
		Run:                     args.Run,
		Connector:               args.Connector,
		RootDir:                 args.RootDir,
		Lineage:                 args.Lineage,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSnapshot:    args.InitialStashSnapshot,
		NoPushHook:              !args.PushHook,
	})
}

func continueRunstate(runState *runstate.RunState, args HandleUnfinishedStateArgs) (bool, error) {
	hasConflicts, err := args.Run.Backend.HasConflicts()
	if err != nil {
		return false, err
	}
	if hasConflicts {
		return false, fmt.Errorf(messages.ContinueUnresolvedConflicts)
	}
	return true, runvm.Execute(runvm.ExecuteArgs{
		RunState:                runState,
		Run:                     args.Run,
		Connector:               args.Connector,
		Lineage:                 args.Lineage,
		RootDir:                 args.RootDir,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSnapshot:    args.InitialStashSnapshot,
		NoPushHook:              !args.PushHook,
	})
}

func discardRunstate(rootDir domain.RepoRootDir) (bool, error) {
	err := persistence.Delete(rootDir)
	return false, err
}

func skipRunstate(runState *runstate.RunState, args HandleUnfinishedStateArgs) (bool, error) {
	skipRunState := runState.CreateSkipRunState()
	return true, runvm.Execute(runvm.ExecuteArgs{
		RunState:                &skipRunState,
		Run:                     args.Run,
		Connector:               args.Connector,
		Lineage:                 args.Lineage,
		RootDir:                 args.RootDir,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSnapshot:    args.InitialStashSnapshot,
		NoPushHook:              !args.PushHook,
	})
}
