package interpreter

import (
	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/config"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/gohacks"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/undo/undoconfig"
	"github.com/git-town/git-town/v15/internal/vm/runstate"
	"github.com/git-town/git-town/v15/internal/vm/shared"
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
		stepName := gohacks.TypeName(nextStep)
		if stepName == "SkipCurrentBranchProgram" {
			args.RunState.SkipCurrentBranchProgram()
			continue
		}
		err := nextStep.Run(shared.RunArgs{
			Backend:                         args.Backend,
			Config:                          args.Config,
			Connector:                       args.Connector,
			DialogTestInputs:                args.DialogTestInputs,
			FinalMessages:                   args.FinalMessages,
			Frontend:                        args.Frontend,
			Git:                             args.Git,
			PrependOpcodes:                  args.RunState.RunProgram.Prepend,
			RegisterUndoablePerennialCommit: args.RunState.RegisterUndoablePerennialCommit,
			UpdateInitialBranchLocalSHA:     args.InitialBranchesSnapshot.Branches.UpdateLocalSHA,
		})
		if err != nil {
			return errored(nextStep, err, args)
		}
	}
}

type ExecuteArgs struct {
	Backend                 gitdomain.RunnerQuerier
	CommandsCounter         Mutable[gohacks.Counter]
	Config                  config.ValidatedConfig
	Connector               Option[hostingdomain.Connector]
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
