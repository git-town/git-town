package runvm

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/git-town/git-town/v9/src/vm/runstate"
)

// Execute runs the commands in the given runstate.
func Execute(args ExecuteArgs) error {
	for {
		nextStep := args.RunState.RunSteps.Pop()
		if nextStep == nil {
			return finished(args)
		}
		stepName := gohacks.TypeName(nextStep)
		if stepName == "SkipCurrentBranchSteps" {
			args.RunState.SkipCurrentBranchSteps()
			continue
		}
		err := nextStep.Run(step.RunArgs{
			AddSteps:                        args.RunState.RunSteps.Prepend,
			Runner:                          args.Run,
			Connector:                       args.Connector,
			Lineage:                         args.Lineage,
			RegisterUndoablePerennialCommit: args.RunState.RegisterUndoablePerennialCommit,
			RemoveBranchFromLineage:         args.Lineage.RemoveBranch,
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
	Debug                   bool
	RootDir                 domain.RepoRootDir
	InitialBranchesSnapshot domain.BranchesSnapshot
	InitialConfigSnapshot   undo.ConfigSnapshot
	InitialStashSnapshot    domain.StashSnapshot
	Lineage                 config.Lineage
	NoPushHook              bool
}
