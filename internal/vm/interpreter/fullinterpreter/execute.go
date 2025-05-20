package fullinterpreter

import (
	"github.com/git-town/git-town/v20/internal/cli/dialog/components"
	"github.com/git-town/git-town/v20/internal/config"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/gohacks"
	"github.com/git-town/git-town/v20/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v20/internal/undo/undoconfig"
	"github.com/git-town/git-town/v20/internal/vm/opcodes"
	"github.com/git-town/git-town/v20/internal/vm/runstate"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// Execute runs the commands in the given runstate.
func Execute(args ExecuteArgs) error {
	for {
		nextStep := args.RunState.RunProgram.Pop()
		if nextStep == nil {
			return finished(finishedArgs{
				Backend:         args.Backend,
				CommandsCounter: args.CommandsCounter,
				FinalMessages:   args.FinalMessages,
				Git:             args.Git,
				RootDir:         args.RootDir,
				RunState:        args.RunState,
				Verbose:         args.Verbose,
			})
		}
		if _, ok := nextStep.(*opcodes.ExitToShell); ok {
			return exitToShell(args)
		}
		err := nextStep.Run(shared.RunArgs{
			Backend:                         args.Backend,
			BranchInfos:                     Some(args.InitialBranchesSnapshot.Branches),
			Config:                          NewMutable(&args.Config),
			Connector:                       args.Connector,
			Detached:                        args.Detached,
			DialogTestInputs:                args.DialogTestInputs,
			FinalMessages:                   args.FinalMessages,
			Frontend:                        args.Frontend,
			Git:                             args.Git,
			PrependOpcodes:                  args.RunState.RunProgram.Prepend,
			RegisterUndoablePerennialCommit: args.RunState.RegisterUndoablePerennialCommit,
			UpdateInitialSnapshotLocalSHA:   args.InitialBranchesSnapshot.Branches.UpdateLocalSHA,
		})
		if err != nil {
			return errored(nextStep, err, args)
		}
		args.RunState.UndoAPIProgram = append(args.RunState.UndoAPIProgram, nextStep.UndoExternalChangesProgram()...)
	}
}

type ExecuteArgs struct {
	Backend                 gitdomain.RunnerQuerier
	CommandsCounter         Mutable[gohacks.Counter]
	Config                  config.ValidatedConfig
	Connector               Option[forgedomain.Connector]
	Detached                configdomain.Detached
	DialogTestInputs        components.TestInputs
	FinalMessages           stringslice.Collector
	Frontend                gitdomain.Runner
	Git                     git.Commands
	HasOpenChanges          bool
	InitialBranch           gitdomain.LocalBranchName
	InitialBranchesSnapshot gitdomain.BranchesSnapshot
	InitialConfigSnapshot   undoconfig.ConfigSnapshot
	InitialStashSize        gitdomain.StashSize
	RootDir                 gitdomain.RepoRootDir
	RunState                runstate.RunState
	Verbose                 configdomain.Verbose
}
