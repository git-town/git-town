package runvm

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/undo"
)

// Execute runs the commands in the given runstate.
func Execute(args ExecuteArgs) error {
	for {
		step := args.RunState.RunSteps.Pop()
		if step == nil {
			return finished(args)
		}
		// TODO: remove this once git skip is sunset
		stepName := gohacks.TypeName(step)
		if stepName == "SkipCurrentBranchSteps" {
			args.RunState.SkipCurrentBranchSteps()
			continue
		}
		err := step.Run(step.RunArgs{
			Runner:                          args.Run,
			Connector:                       args.Connector,
			Lineage:                         args.Lineage,
			RegisterUndoablePerennialCommit: args.RunState.RegisterUndoablePerennialCommit,
			UpdateInitialBranchLocalSHA:     args.InitialBranchesSnapshot.Branches.UpdateLocalSHA,
		})
		if err != nil {
			return errored(step, err, args)
		}
	}
}

type ExecuteArgs struct {
	RunState                *runstate.RunState
	Run                     *git.ProdRunner
	Connector               hosting.Connector
	RootDir                 domain.RepoRootDir
	InitialBranchesSnapshot domain.BranchesSnapshot
	InitialConfigSnapshot   undo.ConfigSnapshot
	InitialStashSnapshot    domain.StashSnapshot
	Lineage                 config.Lineage
	NoPushHook              bool
}
