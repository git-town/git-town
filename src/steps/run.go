package steps

import (
	"fmt"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// Run runs the Git Town command described by the given state.
// nolint: gocyclo, gocognit, nestif, funlen
func Run(runState *RunState, repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	for {
		step := runState.RunStepList.Pop()
		if step == nil {
			runState.MarkAsFinished()
			if runState.IsAbort || runState.isUndo {
				err := DeletePreviousRunState(repo)
				if err != nil {
					return fmt.Errorf("cannot delete previous run state: %w", err)
				}
			} else {
				err := SaveRunState(runState, repo)
				if err != nil {
					return fmt.Errorf("cannot save run state: %w", err)
				}
			}
			fmt.Println()
			return nil
		}
		if getTypeName(step) == "*SkipCurrentBranchSteps" {
			runState.SkipCurrentBranchSteps()
			continue
		}
		if getTypeName(step) == "*PushBranchAfterCurrentBranchSteps" {
			err := runState.AddPushBranchStepAfterCurrentBranchSteps(repo)
			if err != nil {
				return err
			}
			continue
		}
		runErr := step.Run(repo, driver)
		if runErr != nil {
			runState.AbortStepList.Append(step.CreateAbortStep())
			if step.ShouldAutomaticallyAbortOnError() {
				cli.PrintError(fmt.Errorf(runErr.Error() + "\nAuto-aborting..."))
				abortRunState := runState.CreateAbortRunState()
				err := Run(&abortRunState, repo, driver)
				if err != nil {
					return fmt.Errorf("cannot run the abort steps: %w", err)
				}
				cli.Exit(step.GetAutomaticAbortError())
			} else {
				runState.RunStepList.Prepend(step.CreateContinueStep())
				err := runState.MarkAsUnfinished(repo)
				if err != nil {
					return err
				}
				currentBranch, err := repo.Silent.CurrentBranch()
				if err != nil {
					return err
				}
				rebasing, err := repo.Silent.HasRebaseInProgress()
				if err != nil {
					return err
				}
				if runState.Command == "sync" && !(rebasing && repo.Config.IsMainBranch(currentBranch)) {
					runState.UnfinishedDetails.CanSkip = true
				}
				err = SaveRunState(runState, repo)
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
		}
		undoStep, err := step.CreateUndoStep(repo)
		if err != nil {
			return fmt.Errorf("cannot create undo step for %q: %w", step, err)
		}
		runState.UndoStepList.Prepend(undoStep)
	}
}
