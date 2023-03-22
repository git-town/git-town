package runstate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// Execute runs the commands in the given runstate.
//
//nolint:nestif
func Execute(runState *RunState, repo *git.ProdRepo, connector hosting.Connector) error {
	for {
		step := runState.RunStepList.Pop()
		if step == nil {
			runState.MarkAsFinished()
			if runState.IsAbort || runState.isUndo {
				err := Delete(&repo.Backend)
				if err != nil {
					return fmt.Errorf("cannot delete previous run state: %w", err)
				}
			} else {
				err := Save(runState, &repo.Backend)
				if err != nil {
					return fmt.Errorf("cannot save run state: %w", err)
				}
			}
			fmt.Println()
			return nil
		}
		if typeName(step) == "*SkipCurrentBranchSteps" {
			runState.SkipCurrentBranchSteps()
			continue
		}
		if typeName(step) == "*PushBranchAfterCurrentBranchSteps" {
			err := runState.AddPushBranchStepAfterCurrentBranchSteps(&repo.Backend)
			if err != nil {
				return err
			}
			continue
		}
		runErr := step.Run(repo, connector)
		if runErr != nil {
			runState.AbortStepList.Append(step.CreateAbortStep())
			if step.ShouldAutomaticallyAbortOnError() {
				cli.PrintError(fmt.Errorf(runErr.Error() + "\nAuto-aborting..."))
				abortRunState := runState.CreateAbortRunState()
				err := Execute(&abortRunState, repo, connector)
				if err != nil {
					return fmt.Errorf("cannot run the abort steps: %w", err)
				}
				return step.CreateAutomaticAbortError()
			}
			runState.RunStepList.Prepend(step.CreateContinueStep())
			err := runState.MarkAsUnfinished(&repo.Backend)
			if err != nil {
				return err
			}
			currentBranch, err := repo.Backend.CurrentBranch()
			if err != nil {
				return err
			}
			rebasing, err := repo.Backend.HasRebaseInProgress()
			if err != nil {
				return err
			}
			if runState.Command == "sync" && !(rebasing && repo.Config.IsMainBranch(currentBranch)) {
				runState.UnfinishedDetails.CanSkip = true
			}
			err = Save(runState, &repo.Backend)
			if err != nil {
				return fmt.Errorf("cannot save run state: %w", err)
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
		undoStep, err := step.CreateUndoStep(&repo.Backend)
		if err != nil {
			return fmt.Errorf("cannot create undo step for %q: %w", step, err)
		}
		runState.UndoStepList.Prepend(undoStep)
	}
}
