package runstate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/steps"
)

// Execute runs the commands in the given runstate.
func Execute(args ExecuteArgs) error {
	for {
		step := args.RunState.RunStepList.Pop()
		if step == nil {
			return finished(args)
		}
		stepName := typeName(step)
		if stepName == "*SkipCurrentBranchSteps" {
			args.RunState.SkipCurrentBranchSteps()
			continue
		}
		if stepName == "*PushBranchAfterCurrentBranchSteps" {
			err := args.RunState.AddPushBranchStepAfterCurrentBranchSteps(&args.Run.Backend)
			if err != nil {
				return err
			}
			continue
		}
		err := step.Run(args.Run, args.Connector)
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

// finished is called when executing all steps has successfully finished.
func finished(args ExecuteArgs) error {
	args.RunState.MarkAsFinished()
	if args.RunState.IsAbort || args.RunState.isUndo {
		err := Delete(args.RootDir)
		if err != nil {
			return fmt.Errorf(messages.RunstateDeleteProblem, err)
		}
	} else {
		err := Save(args.RunState, args.RootDir)
		if err != nil {
			return fmt.Errorf(messages.RunstateSaveProblem, err)
		}
	}
	fmt.Println()
	args.Run.Stats.PrintAnalysis()
	return nil
}

// errored is called when the given step has resulted in the given error.
func errored(step steps.Step, runErr error, args ExecuteArgs) error {
	args.RunState.AbortStepList.Append(step.CreateAbortSteps()...)
	if step.ShouldAutomaticallyAbortOnError() {
		return autoAbort(step, runErr, args)
	}
	args.RunState.RunStepList.Prepend(step.CreateContinueSteps()...)
	err := args.RunState.MarkAsUnfinished(&args.Run.Backend)
	if err != nil {
		return err
	}
	currentBranch, err := args.Run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	rebasing, err := args.Run.Backend.HasRebaseInProgress()
	if err != nil {
		return err
	}
	if args.RunState.Command == "sync" && !(rebasing && args.Run.Config.IsMainBranch(currentBranch)) {
		args.RunState.UnfinishedDetails.CanSkip = true
	}
	err = Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	message := runErr.Error() + `

To abort, run "git-town abort".
To continue after having resolved conflicts, run "git-town continue".
`
	if args.RunState.UnfinishedDetails.CanSkip {
		message += `To continue by skipping the current branch, run "git-town skip".`
	}
	message += "\n"
	return fmt.Errorf(message)
}

// autoAbort is called when a step that produced an error triggers an auto-abort.
func autoAbort(step steps.Step, runErr error, args ExecuteArgs) error {
	cli.PrintError(fmt.Errorf(messages.RunAutoAborting, runErr.Error()))
	abortRunState := args.RunState.CreateAbortRunState()
	err := Execute(ExecuteArgs{
		RunState:  &abortRunState,
		Run:       args.Run,
		Connector: args.Connector,
		RootDir:   args.RootDir,
	})
	if err != nil {
		return fmt.Errorf(messages.RunstateAbortStepProblem, err)
	}
	return step.CreateAutomaticAbortError()
}

type ExecuteArgs struct {
	RunState  *RunState
	Run       *git.ProdRunner
	Connector hosting.Connector
	RootDir   string
}
