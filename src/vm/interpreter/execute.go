package interpreter

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/hosting"
	"github.com/git-town/git-town/v11/src/undo"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// Execute runs the commands in the given runstate.
func Execute(args ExecuteArgs) error {
	for {
		nextStep := args.RunState.RunProgram.Pop()
		if nextStep == nil {
			return finished(args)
		}
		stepName := gohacks.TypeName(nextStep)
		if stepName == "SkipCurrentBranchProgram" {
			args.RunState.SkipCurrentBranchProgram()
			continue
		}
		err := nextStep.Run(shared.RunArgs{
			PrependOpcodes:                  args.RunState.RunProgram.Prepend,
			Runner:                          args.Run,
			Connector:                       args.Connector,
			Lineage:                         args.Lineage,
			RegisterUndoablePerennialCommit: args.RunState.RegisterUndoablePerennialCommit,
			UpdateInitialBranchLocalSHA:     args.InitialBranchesSnapshot.Branches.UpdateLocalSHA,
		})
		if err != nil {
			return errored(nextStep, err, args)
		}
	}
}

type ExecuteArgs struct {
	RunState                *runstate.RunState
	Run                     *git.ProdRunner
	Connector               hosting.Connector
	Verbose                 bool
	RootDir                 domain.RepoRootDir
	InitialBranchesSnapshot domain.BranchesSnapshot
	InitialConfigSnapshot   undo.ConfigSnapshot
	InitialStashSnapshot    domain.StashSnapshot
	Lineage                 configdomain.Lineage
	NoPushHook              configdomain.NoPushHook
}
