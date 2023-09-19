package runvm

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
)

// Execute runs the commands in the given runstate.
func Execute(args ExecuteArgs) error {
	for {
		step := args.RunState.RunStepList.Pop()
		if step == nil {
			return finished(args)
		}
		stepName := gohacks.TypeName(step)
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
		undoSteps, err := step.CreateUndoSteps(&args.Run.Backend)
		if err != nil {
			return fmt.Errorf(messages.UndoCreateStepProblem, step, err)
		}
		args.RunState.UndoStepList.Prepend(undoSteps...)
	}
}

type ExecuteArgs struct {
	RunState  *runstate.RunState
	Run       *git.ProdRunner
	Connector hosting.Connector
	RootDir   string
	Lineage   config.Lineage
}
