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
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
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
	unfinishedDetails, hasUnfinishedDetails := runState.UnfinishedDetails.Get()
	if !hasUnfinishedDetails {
		return false, nil
	}
	response, exit, err := dialog.AskHowToHandleUnfinishedRunState(
		runState.Command,
		unfinishedDetails.EndBranch,
		unfinishedDetails.EndTime,
		unfinishedDetails.CanSkip,
		args.DialogTestInputs.Value.Next(),
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
		return undoRunState(args, runState)
	case dialog.ResponseSkip:
		return skipRunstate(args, runState)
	case dialog.ResponseQuit:
		return true, nil
	}
	return false, fmt.Errorf(messages.DialogUnexpectedResponse, response)
}

type UnfinishedStateArgs struct {
	Backend           gitdomain.RunnerQuerier
	CommandsCounter   Mutable[gohacks.Counter]
	Connector         Option[hostingdomain.Connector]
	DialogTestInputs  Mutable[components.TestInputs]
	FinalMessages     stringslice.Collector
	Frontend          gitdomain.Runner
	Git               git.Commands
	HasOpenChanges    bool
	PushHook          configdomain.PushHook
	RepoStatus        gitdomain.RepoStatus
	RootDir           gitdomain.RepoRootDir
	UnvalidatedConfig config.UnvalidatedConfig
	Verbose           configdomain.Verbose
}

func continueRunstate(runState runstate.RunState, args UnfinishedStateArgs) (bool, error) {
	if args.RepoStatus.Conflicts {
		return false, errors.New(messages.ContinueUnresolvedConflicts)
	}
	validatedConfig, exit, err := quickValidateConfig(quickValidateConfigArgs{
		backend:      args.Backend,
		dialogInputs: args.DialogTestInputs,
		git:          args.Git,
		unvalidated:  args.UnvalidatedConfig,
	})
	if err != nil || exit {
		return exit, err
	}
	return true, fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 args.Backend,
		CommandsCounter:         args.CommandsCounter,
		Config:                  validatedConfig,
		Connector:               args.Connector,
		DialogTestInputs:        args.DialogTestInputs,
		FinalMessages:           args.FinalMessages,
		Frontend:                args.Frontend,
		Git:                     args.Git,
		HasOpenChanges:          args.RepoStatus.OpenChanges,
		InitialBranch:           runState.BeginBranchesSnapshot.Active.GetOrPanic(),
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

// quickly provides a ValidatedConfig instance in situations where we continue runstate.
// It is expected that all data exists.
// This doesn't change lineage since we are in the middle of an ongoing Git Town operation.
func quickValidateConfig(args quickValidateConfigArgs) (config.ValidatedConfig, bool, error) {
	mainBranch, hasMain := args.unvalidated.Config.Value.MainBranch.Get()
	if !hasMain {
		branchesSnapshot, err := args.git.BranchesSnapshot(args.backend)
		if err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
		localBranches := branchesSnapshot.Branches.LocalBranches().Names()
		validatedMain, exit, err := dialog.MainBranch(localBranches, args.git.DefaultBranch(args.backend), args.dialogInputs.Value.Next())
		if err != nil || exit {
			return config.EmptyValidatedConfig(), exit, err
		}
		if err = args.unvalidated.SetMainBranch(validatedMain); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
		mainBranch = validatedMain
	}
	gitUserEmail, gitUserName, err := GitUser(args.unvalidated.Config.Get())
	if err != nil {
		return config.EmptyValidatedConfig(), false, err
	}
	return config.ValidatedConfig{
		Config: configdomain.ValidatedConfig{
			UnvalidatedConfig: args.unvalidated.Config.Value,
			GitUserEmail:      gitUserEmail,
			GitUserName:       gitUserName,
			MainBranch:        mainBranch,
		},
		UnvalidatedConfig: &args.unvalidated,
	}, false, nil
}

func skipRunstate(args UnfinishedStateArgs, runState runstate.RunState) (bool, error) {
	currentBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return false, err
	}
	validatedConfig, exit, err := quickValidateConfig(quickValidateConfigArgs{
		backend:      args.Backend,
		dialogInputs: args.DialogTestInputs,
		git:          args.Git,
		unvalidated:  args.UnvalidatedConfig,
	})
	if err != nil || exit {
		return exit, err
	}
	return true, skip.Execute(skip.ExecuteArgs{
		Backend:         args.Backend,
		CommandsCounter: args.CommandsCounter,
		Config:          validatedConfig,
		Connector:       args.Connector,
		FinalMessages:   args.FinalMessages,
		Frontend:        args.Frontend,
		Git:             args.Git,
		HasOpenChanges:  args.HasOpenChanges,
		InitialBranch:   currentBranch,
		RootDir:         args.RootDir,
		RunState:        runState,
		TestInputs:      args.DialogTestInputs,
		Verbose:         args.Verbose,
	})
}

func undoRunState(args UnfinishedStateArgs, runState runstate.RunState) (bool, error) {
	validatedConfig, exit, err := quickValidateConfig(quickValidateConfigArgs{
		backend:      args.Backend,
		dialogInputs: args.DialogTestInputs,
		git:          args.Git,
		unvalidated:  args.UnvalidatedConfig,
	})
	if err != nil || exit {
		return exit, err
	}
	return true, undo.Execute(undo.ExecuteArgs{
		Backend:          args.Backend,
		CommandsCounter:  args.CommandsCounter,
		Config:           validatedConfig,
		FinalMessages:    args.FinalMessages,
		Frontend:         args.Frontend,
		Git:              args.Git,
		HasOpenChanges:   args.HasOpenChanges,
		InitialStashSize: runState.BeginStashSize,
		RootDir:          args.RootDir,
		RunState:         runState,
		Verbose:          args.Verbose,
	})
}

type quickValidateConfigArgs struct {
	backend      gitdomain.RunnerQuerier
	dialogInputs Mutable[components.TestInputs]
	git          git.Commands
	unvalidated  config.UnvalidatedConfig
}
