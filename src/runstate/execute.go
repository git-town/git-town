package runstate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
)

// Execute runs the commands in the given runstate.
//
//nolint:nestif
func Execute(runState *RunState, run *git.ProdRunner, connector hosting.Connector, rootDir string) error {
	for {
		step := runState.RunStepList.Pop()
		if step == nil {
			runState.MarkAsFinished()
			if runState.IsAbort || runState.isUndo {
				err := Delete(rootDir)
				if err != nil {
					return fmt.Errorf(messages.RunstateDeleteProblem, err)
				}
			} else {
				err := Save(runState, rootDir)
				if err != nil {
					return fmt.Errorf(messages.RunstateSaveProblem, err)
				}
			}
			fmt.Println()
			run.Stats.PrintAnalysis()
			return nil
		}
		if typeName(step) == "*SkipCurrentBranchSteps" {
			runState.SkipCurrentBranchSteps()
			continue
		}
		if typeName(step) == "*PushBranchAfterCurrentBranchSteps" {
			err := runState.AddPushBranchStepAfterCurrentBranchSteps(&run.Backend)
			if err != nil {
				return err
			}
			continue
		}
		runErr := step.Run(run, connector)
		if runErr != nil {
			runState.AbortStepList.Append(step.CreateAbortStep())
			if step.ShouldAutomaticallyAbortOnError() {
				cli.PrintError(fmt.Errorf(runErr.Error() + "\nAuto-aborting..."))
				abortRunState := runState.CreateAbortRunState()
				err := Execute(&abortRunState, run, connector, rootDir)
				if err != nil {
					return fmt.Errorf(messages.RunstateAbortStepProblem, err)
				}
				return step.CreateAutomaticAbortError()
			}
			runState.RunStepList.Prepend(step.CreateContinueStep())
			err := runState.MarkAsUnfinished(&run.Backend)
			if err != nil {
				return err
			}
			currentBranch, err := run.Backend.CurrentBranch()
			if err != nil {
				return err
			}
			rebasing, err := run.Backend.HasRebaseInProgress()
			if err != nil {
				return err
			}
			if runState.Command == "sync" && !(rebasing && run.Config.IsMainBranch(currentBranch)) {
				runState.UnfinishedDetails.CanSkip = true
			}
			err = Save(runState, rootDir)
			if err != nil {
				return fmt.Errorf(messages.RunstateSaveProblem, err)
			}
			message := runErr.Error() + `

To abort, run "git-town abort".
To continue after having resolved conflicts, run "git-town continue".
`
			if runState.UnfinishedDetails.CanSkip {
				message += `To continue by skipping the current branch, run "git-town skip".`
			}
			message += "\n"
			return fmt.Errorf(message)
		}
		undoStep, err := step.CreateUndoStep(&run.Backend)
		if err != nil {
			return fmt.Errorf(messages.UndoCreateStepProblem, step, err)
		}
		runState.UndoStepList.Prepend(undoStep)
	}
}
