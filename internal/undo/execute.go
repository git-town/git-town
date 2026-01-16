package undo

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/lightinterpreter"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ExecuteArgs struct {
	Backend          subshelldomain.RunnerQuerier
	CommandsCounter  Mutable[gohacks.Counter]
	Config           config.ValidatedConfig
	ConfigDir        configdomain.RepoConfigDir
	Connector        Option[forgedomain.Connector]
	FinalMessages    stringslice.Collector
	Frontend         subshelldomain.Runner
	Git              git.Commands
	HasOpenChanges   bool
	InitialStashSize gitdomain.StashSize
	RunState         runstate.RunState
}

// undoes the persisted runstate
func Execute(args ExecuteArgs) (changedBranches gitdomain.LocalBranchNames, err error) {
	if args.RunState.DryRun {
		return gitdomain.LocalBranchNames{}, nil
	}
	program, changedBranches := CreateUndoForFinishedProgram(CreateUndoProgramArgs{
		Backend:        args.Backend,
		Config:         args.Config,
		DryRun:         args.RunState.DryRun,
		FinalMessages:  args.FinalMessages,
		Git:            args.Git,
		HasOpenChanges: args.HasOpenChanges,
		PushHook:       args.Config.NormalConfig.PushHook,
		RunState:       args.RunState,
	})
	lightinterpreter.Execute(lightinterpreter.ExecuteArgs{
		Backend:       args.Backend,
		BranchInfos:   args.RunState.BeginBranchesSnapshot.Branches,
		Config:        args.Config,
		Connector:     args.Connector,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Git:           args.Git,
		Prog:          program,
	})
	runstatePath := runstate.NewRunstatePath(args.ConfigDir)
	err := os.Remove(runstatePath.String())
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf(messages.RunstateDeleteProblem, err)
	}
	print.Footer(args.Config.NormalConfig.Verbose, args.CommandsCounter.Immutable(), args.FinalMessages.Result())
	return nil
}
