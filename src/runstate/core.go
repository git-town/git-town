// Package runstate represents the current execution status of a Git Town command and allows to run or persist it to disk.
// This is used by the "abort", "continue", and "undo" commands.
// The central data structure is RunState.
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
	AbortStepList     StepList `exhaustruct:"optional"`
	Command           string
	IsAbort           bool `exhaustruct:"optional"`
	isUndo            bool `exhaustruct:"optional"`
	RunStepList       StepList
	UndoStepList      StepList                   `exhaustruct:"optional"`
	UnfinishedDetails *UnfinishedRunStateDetails `exhaustruct:"optional"`
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
		if !isCheckoutStep(step) {
			popped.Append(runState.RunStepList.Pop())
		} else {
			currentBranch, err := repo.Silent.CurrentBranch()
			if err != nil {
				return err
			}
			runState.RunStepList.Prepend(&steps.PushBranchStep{Branch: currentBranch})
			runState.RunStepList.PrependList(popped)
			break
		}
	}
	return nil
}

// CreateAbortRunState returns a new runstate
// to be run to aborting and undoing the Git Town command
// represented by this runstate.
func (runState *RunState) CreateAbortRunState() RunState {
	stepList := runState.AbortStepList
	stepList.AppendList(runState.UndoStepList)
	return RunState{
		Command:     runState.Command,
		IsAbort:     true,
		RunStepList: stepList,
	}
}

// CreateSkipRunState returns a new Runstate
// that skips operations for the current branch.
func (runState *RunState) CreateSkipRunState() RunState {
	result := RunState{
		Command:     runState.Command,
		RunStepList: runState.AbortStepList,
	}
	for _, step := range runState.UndoStepList.List {
		if isCheckoutStep(step) {
			break
		}
		result.RunStepList.Append(step)
	}
	skipping := true
	for _, step := range runState.RunStepList.List {
		if isCheckoutStep(step) {
			skipping = false
		}
		if !skipping {
			result.RunStepList.Append(step)
		}
	}
	return result
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

func (runState *RunState) HasAbortSteps() bool {
	return !runState.AbortStepList.IsEmpty()
}

func (runState *RunState) HasRunSteps() bool {
	return !runState.RunStepList.IsEmpty()
}

func (runState *RunState) HasUndoSteps() bool {
	return !runState.UndoStepList.IsEmpty()
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
		if isCheckoutStep(step) {
			break
		}
		runState.RunStepList.Pop()
	}
}

func isCheckoutStep(step steps.Step) bool {
	return typeName(step) == "*CheckoutStep"
}
