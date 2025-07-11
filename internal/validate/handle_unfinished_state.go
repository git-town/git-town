package validate

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/skip"
	"github.com/git-town/git-town/v21/internal/state"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/undo"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/fullinterpreter"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
func HandleUnfinishedState(args UnfinishedStateArgs) (dialogdomain.Exit, error) {
	runState, hasRunState := args.RunState.Get()
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
		args.DialogTestInputs.Next(),
	)
	if err != nil {
		return false, err
	}
	if exit {
		return exit, errors.New("user aborted")
	}
	// Create the connector now if the Git Town command hasn't provided one yet.
	if args.Connector.IsNone() {
		normalConfig := args.UnvalidatedConfig.NormalConfig
		args.Connector, err = forge.NewConnector(forge.NewConnectorArgs{
			Backend:              args.Backend,
			BitbucketAppPassword: normalConfig.BitbucketAppPassword,
			BitbucketUsername:    normalConfig.BitbucketUsername,
			CodebergToken:        normalConfig.CodebergToken,
			ForgeType:            normalConfig.ForgeType,
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
			return false, err
		}
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
	Backend           subshelldomain.RunnerQuerier
	CommandsCounter   Mutable[gohacks.Counter]
	Connector         Option[forgedomain.Connector]
	Detached          configdomain.Detached
	DialogTestInputs  dialogcomponents.TestInputs
	FinalMessages     stringslice.Collector
	Frontend          subshelldomain.Runner
	Git               git.Commands
	HasOpenChanges    bool
	PushHook          configdomain.PushHook
	RepoStatus        gitdomain.RepoStatus
	RootDir           gitdomain.RepoRootDir
	RunState          Option[runstate.RunState]
	UnvalidatedConfig config.UnvalidatedConfig
	Verbose           configdomain.Verbose
}

func continueRunstate(runState runstate.RunState, args UnfinishedStateArgs) (dialogdomain.Exit, error) {
	if args.RepoStatus.Conflicts {
		return false, errors.New(messages.ContinueUnresolvedConflicts)
	}
	validatedConfig, exit, err := quickValidateConfig(quickValidateConfigArgs{
		backend:      args.Backend,
		dialogInputs: args.DialogTestInputs,
		git:          args.Git,
		unvalidated:  NewMutable(&args.UnvalidatedConfig),
	})
	if err != nil || exit {
		return exit, err
	}
	return true, fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 args.Backend,
		CommandsCounter:         args.CommandsCounter,
		Config:                  validatedConfig,
		Connector:               args.Connector,
		Detached:                args.Detached,
		DialogTestInputs:        args.DialogTestInputs,
		FinalMessages:           args.FinalMessages,
		Frontend:                args.Frontend,
		Git:                     args.Git,
		HasOpenChanges:          args.RepoStatus.OpenChanges,
		InitialBranch:           runState.BeginBranchesSnapshot.Active.GetOrPanic(),
		InitialBranchesSnapshot: runState.BeginBranchesSnapshot,
		InitialConfigSnapshot:   runState.BeginConfigSnapshot,
		InitialStashSize:        runState.BeginStashSize,
		PendingCommand:          Some(runState.Command),
		RootDir:                 args.RootDir,
		RunState:                runState,
		Verbose:                 args.Verbose,
	})
}

func discardRunstate(rootDir gitdomain.RepoRootDir) (dialogdomain.Exit, error) {
	_, err := state.Delete(rootDir, state.FileTypeRunstate)
	return false, err
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
		localBranches := branchesSnapshot.Branches.LocalBranches().Names()
		validatedMain, exit, err := dialog.MainBranch(localBranches, gitconfig.DefaultBranch(args.backend), args.dialogInputs.Next())
		if err != nil || exit {
			return config.EmptyValidatedConfig(), exit, err
		}
		if err = args.unvalidated.Value.SetMainBranch(validatedMain, args.backend); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
		mainBranch = validatedMain
	}
	gitUserEmail, gitUserName, err := GitUser(args.unvalidated.Value.UnvalidatedConfig)
	if err != nil {
		return config.EmptyValidatedConfig(), false, err
	}
	return config.ValidatedConfig{
		ValidatedConfigData: configdomain.ValidatedConfigData{
			GitUserEmail: gitUserEmail,
			GitUserName:  gitUserName,
			MainBranch:   mainBranch,
		},
		NormalConfig: args.unvalidated.Value.NormalConfig,
	}, false, nil
}

func skipRunstate(args UnfinishedStateArgs, runState runstate.RunState) (dialogdomain.Exit, error) {
	currentBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return false, err
	}
	validatedConfig, exit, err := quickValidateConfig(quickValidateConfigArgs{
		backend:      args.Backend,
		dialogInputs: args.DialogTestInputs,
		git:          args.Git,
		unvalidated:  NewMutable(&args.UnvalidatedConfig),
	})
	if err != nil || exit {
		return exit, err
	}
	return true, skip.Execute(skip.ExecuteArgs{
		Backend:         args.Backend,
		CommandsCounter: args.CommandsCounter,
		Config:          validatedConfig,
		Connector:       args.Connector,
		Detached:        args.Detached,
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

func undoRunState(args UnfinishedStateArgs, runState runstate.RunState) (dialogdomain.Exit, error) {
	validatedConfig, exit, err := quickValidateConfig(quickValidateConfigArgs{
		backend:      args.Backend,
		dialogInputs: args.DialogTestInputs,
		git:          args.Git,
		unvalidated:  NewMutable(&args.UnvalidatedConfig),
	})
	if err != nil || exit {
		return exit, err
	}
	return true, undo.Execute(undo.ExecuteArgs{
		Backend:          args.Backend,
		CommandsCounter:  args.CommandsCounter,
		Config:           validatedConfig,
		Connector:        args.Connector,
		Detached:         args.Detached,
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
	backend      subshelldomain.RunnerQuerier
	dialogInputs dialogcomponents.TestInputs
	git          git.Commands
	unvalidated  Mutable[config.UnvalidatedConfig]
}
