package runvm

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/undo"
)

// Execute runs the commands in the given runstate.
func Execute(args ExecuteArgs) error {
	for {
		step := args.RunState.RunStepList.Pop()
		if step == nil {
			return finished(args)
		}
		stepName := runstate.TypeName(step)
		if stepName == "SkipCurrentBranchSteps" {
			args.RunState.SkipCurrentBranchSteps()
			continue
		}
		if stepName == "PushBranchAfterCurrentBranchSteps" {
			err := args.RunState.AddPushBranchStepAfterCurrentBranchSteps(&args.Run.Backend)
			if err != nil {
				return err
			}
			continue
		}
		err := step.Run(steps.RunArgs{
			Runner:    args.Run,
			Connector: args.Connector,
			Lineage:   args.Lineage,
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
	RootDir                 string
	InitialBranchesSnapshot undo.BranchesSnapshot
	InitialConfigSnapshot   undo.ConfigSnapshot
	Lineage                 config.Lineage
}
