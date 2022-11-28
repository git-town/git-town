package runstate

import (
	"time"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/steps"
)

// UnfinishedRunStateDetails has details about an unfinished run state.
type UnfinishedRunStateDetails struct {
	CanSkip   bool
	EndBranch string
	EndTime   time.Time
}

// RunState represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has been done so far.
type RunState struct {
	AbortStepList     StepList
	Command           string
	IsAbort           bool
	isUndo            bool
	RunStepList       StepList
	UndoStepList      StepList
	UnfinishedDetails *UnfinishedRunStateDetails
}

// New constructs a RunState instance with the given values.
func New(command string, stepList StepList) *RunState {
	return &RunState{
		Command:     command,
		RunStepList: stepList,
	}
}

// AddPushBranchStepAfterCurrentBranchSteps inserts a PushBranchStep
// after all the steps for the current branch.
func (runState *RunState) AddPushBranchStepAfterCurrentBranchSteps(repo *git.ProdRepo) error {
	popped := StepList{}
	for {
		step := runState.RunStepList.Peek()
		if !isCheckoutBranchStep(step) {
			popped.Append(runState.RunStepList.Pop())
		} else {
			currentBranch, err := repo.Silent.CurrentBranch()
			if err != nil {
				return err
			}
			runState.RunStepList.Prepend(&steps.PushBranchStep{BranchName: currentBranch})
			runState.RunStepList.PrependList(popped)
			break
		}
	}
	return nil
}

// CreateAbortRunState returns a new runstate
// to be run to aborting and undoing the Git Town command
// represented by this runstate.
func (runState *RunState) CreateAbortRunState() (result RunState) {
	result.Command = runState.Command
	result.IsAbort = true
	result.RunStepList.AppendList(runState.AbortStepList)
	result.RunStepList.AppendList(runState.UndoStepList)
	return
}

// CreateSkipRunState returns a new Runstate
// that skips operations for the current branch.
func (runState *RunState) CreateSkipRunState() (result RunState) {
	result.Command = runState.Command
	result.RunStepList.AppendList(runState.AbortStepList)
	for _, step := range runState.UndoStepList.List {
		if isCheckoutBranchStep(step) {
			break
		}
		result.RunStepList.Append(step)
	}
	skipping := true
	for _, step := range runState.RunStepList.List {
		if isCheckoutBranchStep(step) {
			skipping = false
		}
		if !skipping {
			result.RunStepList.Append(step)
		}
	}
	return
}

// CreateUndoRunState returns a new runstate
// to be run when undoing the Git Town command
// represented by this runstate.
func (runState *RunState) CreateUndoRunState() RunState {
	return RunState{
		Command:     runState.Command,
		isUndo:      true,
		RunStepList: runState.UndoStepList,
	}
}

// IsUnfinished returns whether or not the run state is unfinished.
func (runState *RunState) IsUnfinished() bool {
	return runState.UnfinishedDetails != nil
}

// MarkAsFinished updates the run state to be marked as finished.
func (runState *RunState) MarkAsFinished() {
	runState.UnfinishedDetails = nil
}

// MarkAsUnfinished updates the run state to be marked as unfinished and populates informational fields.
func (runState *RunState) MarkAsUnfinished(repo *git.ProdRepo) error {
	currentBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return err
	}
	runState.UnfinishedDetails = &UnfinishedRunStateDetails{
		CanSkip:   false,
		EndBranch: currentBranch,
		EndTime:   time.Now(),
	}
	return nil
}

// SkipCurrentBranchSteps removes the steps for the current branch
// from this run state.
func (runState *RunState) SkipCurrentBranchSteps() {
	for {
		step := runState.RunStepList.Peek()
		if isCheckoutBranchStep(step) {
			break
		}
		runState.RunStepList.Pop()
	}
}

func isCheckoutBranchStep(step steps.Step) bool {
	return typeName(step) == "*CheckoutBranchStep"
}
