package validate

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/skip"
	"github.com/git-town/git-town/v22/internal/state"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/internal/undo"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/fullinterpreter"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
func HandleUnfinishedState(args UnfinishedStateArgs) (configdomain.ProgramFlow, error) {
	runState, hasRunState := args.RunState.Get()
	if !hasRunState || runState.IsFinished() {
		return configdomain.ProgramFlowContinue, nil
	}
	unfinishedDetails, hasUnfinishedDetails := runState.UnfinishedDetails.Get()
	if !hasUnfinishedDetails {
		return configdomain.ProgramFlowContinue, nil
	}
	response, exit, err := dialog.AskHowToHandleUnfinishedRunState(
		runState.Command,
		unfinishedDetails.EndBranch,
		unfinishedDetails.EndTime,
		unfinishedDetails.CanSkip,
		args.Inputs,
	)
	if err != nil {
		return configdomain.ProgramFlowExit, err
	}
	if exit {
		return configdomain.ProgramFlowExit, errors.New("user aborted")
	}
	// Create the connector now if the Git Town command hasn't provided one yet.
	if args.Connector.IsNone() {
		normalConfig := args.UnvalidatedConfig.NormalConfig
		args.Connector, err = forge.NewConnector(forge.NewConnectorArgs{
			Backend:              args.Backend,
			BitbucketAppPassword: normalConfig.BitbucketAppPassword,
			BitbucketUsername:    normalConfig.BitbucketUsername,
			ForgeType:            normalConfig.ForgeType,
			ForgejoToken:         normalConfig.ForgejoToken,
			Frontend:             args.Frontend,
			GitHubConnectorType:  normalConfig.GitHubConnectorType,
			GitHubToken:          normalConfig.GitHubToken,
			GitLabConnectorType:  normalConfig.GitLabConnectorType,
			GitLabToken:          normalConfig.GitLabToken,
			GiteaToken:           normalConfig.GiteaToken,
			Log:                  print.Logger{},
			RemoteURL:            normalConfig.DevURL(args.Backend),
		})
		if err != nil {
			return configdomain.ProgramFlowExit, err
		}
	}
	switch response {
	case dialog.ResponseBoth:
		_, err := continueRunstate(runState, args)
		return configdomain.ProgramFlowContinue, err
	case dialog.ResponseDiscard:
		return discardRunstate(args.RootDir)
	case dialog.ResponseContinue:
		return continueRunstate(runState, args)
	case dialog.ResponseUndo:
		return undoRunState(args, runState)
	case dialog.ResponseSkip:
		return skipRunstate(args, runState)
	case dialog.ResponseQuit:
		return configdomain.ProgramFlowExit, nil
	}
	return configdomain.ProgramFlowExit, fmt.Errorf(messages.DialogUnexpectedResponse, response)
}

type UnfinishedStateArgs struct {
	Backend           subshelldomain.RunnerQuerier
	CommandsCounter   Mutable[gohacks.Counter]
	Connector         Option[forgedomain.Connector]
	FinalMessages     stringslice.Collector
	Frontend          subshelldomain.Runner
	Git               git.Commands
	HasOpenChanges    bool
	Inputs            dialogcomponents.Inputs
	PushHook          configdomain.PushHook
	RepoStatus        gitdomain.RepoStatus
	RootDir           gitdomain.RepoRootDir
	RunState          Option[runstate.RunState]
	UnvalidatedConfig config.UnvalidatedConfig
}

func continueRunstate(runState runstate.RunState, args UnfinishedStateArgs) (configdomain.ProgramFlow, error) {
	if args.RepoStatus.Conflicts {
		return configdomain.ProgramFlowExit, errors.New(messages.ContinueUnresolvedConflicts)
	}
	validatedConfig, exit, err := quickValidateConfig(quickValidateConfigArgs{
		backend:     args.Backend,
		git:         args.Git,
		inputs:      args.Inputs,
		unvalidated: NewMutable(&args.UnvalidatedConfig),
	})
	if err != nil || exit {
		return configdomain.ProgramFlowExit, err
	}
	return configdomain.ProgramFlowExit, fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 args.Backend,
		CommandsCounter:         args.CommandsCounter,
		Config:                  validatedConfig,
		Connector:               args.Connector,
		FinalMessages:           args.FinalMessages,
		Frontend:                args.Frontend,
		Git:                     args.Git,
		HasOpenChanges:          args.RepoStatus.OpenChanges,
		InitialBranch:           runState.BeginBranchesSnapshot.Active.GetOrPanic(),
		InitialBranchesSnapshot: runState.BeginBranchesSnapshot,
		InitialConfigSnapshot:   runState.BeginConfigSnapshot,
		InitialStashSize:        runState.BeginStashSize,
		Inputs:                  args.Inputs,
		PendingCommand:          Some(runState.Command),
		RootDir:                 args.RootDir,
		RunState:                runState,
	})
}

func discardRunstate(rootDir gitdomain.RepoRootDir) (configdomain.ProgramFlow, error) {
	_, err := state.Delete(rootDir, state.FileTypeRunstate)
	return configdomain.ProgramFlowContinue, err
}

// quickly provides a ValidatedConfig instance in situations where we continue runstate.
// It is expected that all data exists.
// This doesn't change lineage since we are in the middle of an ongoing Git Town operation.
func quickValidateConfig(args quickValidateConfigArgs) (config.ValidatedConfig, dialogdomain.Exit, error) {
	mainBranch, hasMain := args.unvalidated.Value.UnvalidatedConfig.MainBranch.Get()
	if !hasMain {
		branchesSnapshot, err := args.git.BranchesSnapshot(args.backend)
		if err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
		localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
		var exit dialogdomain.Exit
		_, mainBranch, exit, err = dialog.MainBranch(dialog.MainBranchArgs{
			Inputs:         args.inputs,
			Local:          args.unvalidated.Value.GitGlobal.MainBranch,
			LocalBranches:  localBranches,
			StandardBranch: args.git.StandardBranch(args.backend),
			Unscoped:       args.unvalidated.Value.GitUnscoped.MainBranch,
		})
		if err != nil || exit {
			return config.EmptyValidatedConfig(), exit, err
		}
		if err = args.unvalidated.Value.SetMainBranch(mainBranch, args.backend); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
	}
	return config.ValidatedConfig{
		ValidatedConfigData: configdomain.ValidatedConfigData{
			MainBranch: mainBranch,
		},
		NormalConfig: args.unvalidated.Value.NormalConfig,
	}, false, nil
}

func skipRunstate(args UnfinishedStateArgs, runState runstate.RunState) (configdomain.ProgramFlow, error) {
	currentBranchOpt, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return configdomain.ProgramFlowExit, err
	}
	currentBranch, hasCurrentBranch := currentBranchOpt.Get()
	if !hasCurrentBranch {
		return configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	validatedConfig, exit, err := quickValidateConfig(quickValidateConfigArgs{
		backend:     args.Backend,
		git:         args.Git,
		inputs:      args.Inputs,
		unvalidated: NewMutable(&args.UnvalidatedConfig),
	})
	if err != nil || exit {
		return configdomain.ProgramFlowExit, err
	}
	return configdomain.ProgramFlowExit, skip.Execute(skip.ExecuteArgs{
		Backend:         args.Backend,
		CommandsCounter: args.CommandsCounter,
		Config:          validatedConfig,
		Connector:       args.Connector,
		FinalMessages:   args.FinalMessages,
		Frontend:        args.Frontend,
		Git:             args.Git,
		HasOpenChanges:  args.HasOpenChanges,
		InitialBranch:   currentBranch,
		Inputs:          args.Inputs,
		RootDir:         args.RootDir,
		RunState:        runState,
	})
}

func undoRunState(args UnfinishedStateArgs, runState runstate.RunState) (configdomain.ProgramFlow, error) {
	validatedConfig, exit, err := quickValidateConfig(quickValidateConfigArgs{
		backend:     args.Backend,
		git:         args.Git,
		inputs:      args.Inputs,
		unvalidated: NewMutable(&args.UnvalidatedConfig),
	})
	if err != nil || exit {
		return configdomain.ProgramFlowExit, err
	}
	return configdomain.ProgramFlowExit, undo.Execute(undo.ExecuteArgs{
		Backend:          args.Backend,
		CommandsCounter:  args.CommandsCounter,
		Config:           validatedConfig,
		Connector:        args.Connector,
		FinalMessages:    args.FinalMessages,
		Frontend:         args.Frontend,
		Git:              args.Git,
		HasOpenChanges:   args.HasOpenChanges,
		InitialStashSize: runState.BeginStashSize,
		RootDir:          args.RootDir,
		RunState:         runState,
	})
}

type quickValidateConfigArgs struct {
	backend     subshelldomain.RunnerQuerier
	git         git.Commands
	inputs      dialogcomponents.Inputs
	unvalidated Mutable[config.UnvalidatedConfig]
}
