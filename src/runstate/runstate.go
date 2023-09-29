package runstate

import (
	"fmt"
	"strings"
	"time"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/slice"
	"github.com/git-town/git-town/v9/src/steps"
)

// RunState represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has been done so far.
// TODO: rename the "XXXStepList" fields to "XXXSteps".
type RunState struct {
	Command                  string                     `json:"Command"`
	IsAbort                  bool                       `exhaustruct:"optional"     json:"IsAbort"`
	IsUndo                   bool                       `exhaustruct:"optional"     json:"IsUndo"`
	AbortStepList            StepList                   `exhaustruct:"optional"     json:"AbortStepList"`
	RunStepList              StepList                   `json:"RunStepList"`
	UndoStepList             StepList                   `exhaustruct:"optional"     json:"UndoStepList"`
	InitialActiveBranch      domain.LocalBranchName     `json:"InitialActiveBranch"`
	FinalUndoStepList        StepList                   `exhaustruct:"optional"     json:"FinalUndoStepList"`
	UnfinishedDetails        *UnfinishedRunStateDetails `exhaustruct:"optional"     json:"UnfinishedDetails"`
	UndoablePerennialCommits []domain.SHA               `exhaustruct:"optional"     json:"UndoablePerennialCommits"`
}

// AddPushBranchStepAfterCurrentBranchSteps inserts a PushBranchStep
// after all the steps for the current branch.
func (runState *RunState) AddPushBranchStepAfterCurrentBranchSteps(backend *git.BackendCommands) error {
	popped := StepList{}
	for {
		step := runState.RunStepList.Peek()
		if !isCheckoutStep(step) {
			popped.Append(runState.RunStepList.Pop())
		} else {
			currentBranch, err := backend.CurrentBranch()
			if err != nil {
				return err
			}
			runState.RunStepList.Prepend(&steps.PushCurrentBranchStep{CurrentBranch: currentBranch, NoPushHook: false})
			runState.RunStepList.PrependList(popped)
			break
		}
	}
	return nil
}

// RegisterUndoablePerennialCommit stores the given commit on a perennial branch as undoable.
// This method is used as a callback.
// TODO: rename runState to rs.
func (runState *RunState) RegisterUndoablePerennialCommit(commit domain.SHA) {
	runState.UndoablePerennialCommits = append(runState.UndoablePerennialCommits, commit)
}

// CreateAbortRunState returns a new runstate
// to be run to aborting and undoing the Git Town command
// represented by this runstate.
func (runState *RunState) CreateAbortRunState() RunState {
	stepList := runState.AbortStepList
	stepList.AppendList(runState.UndoStepList)
	return RunState{
		Command:             runState.Command,
		IsAbort:             true,
		InitialActiveBranch: runState.InitialActiveBranch,
		RunStepList:         stepList,
	}
}

// CreateSkipRunState returns a new Runstate
// that skips operations for the current branch.
func (runState *RunState) CreateSkipRunState() RunState {
	result := RunState{
		Command:             runState.Command,
		InitialActiveBranch: runState.InitialActiveBranch,
		RunStepList:         runState.AbortStepList,
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
	result.RunStepList.List = slice.LowerAll[steps.Step](result.RunStepList.List, &steps.RestoreOpenChangesStep{})
	return result
}

// CreateUndoRunState returns a new runstate
// to be run when undoing the Git Town command
// represented by this runstate.
func (runState *RunState) CreateUndoRunState() RunState {
	result := RunState{
		Command:                  runState.Command,
		InitialActiveBranch:      runState.InitialActiveBranch,
		IsUndo:                   true,
		RunStepList:              runState.UndoStepList,
		UndoablePerennialCommits: []domain.SHA{},
	}
	result.RunStepList.Append(&steps.CheckoutStep{Branch: runState.InitialActiveBranch})
	result.RunStepList = result.RunStepList.RemoveDuplicateCheckoutSteps()
	return result
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
func (runState *RunState) MarkAsUnfinished(backend *git.BackendCommands) error {
	currentBranch, err := backend.CurrentBranch()
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

func (runState *RunState) String() string {
	result := strings.Builder{}
	result.WriteString("RunState:\n")
	result.WriteString("  Command: ")
	result.WriteString(runState.Command)
	result.WriteString("\n  IsAbort: ")
	result.WriteString(fmt.Sprintf("%t", runState.IsAbort))
	result.WriteString("\n  IsUndo: ")
	result.WriteString(fmt.Sprintf("%t", runState.IsUndo))
	result.WriteString("\n  AbortStepList: ")
	result.WriteString(runState.AbortStepList.StringIndented("    "))
	result.WriteString("  RunStepList: ")
	result.WriteString(runState.RunStepList.StringIndented("    "))
	result.WriteString("  UndoStepList: ")
	result.WriteString(runState.UndoStepList.StringIndented("    "))
	if runState.UnfinishedDetails != nil {
		result.WriteString("  UnfineshedDetails: ")
		result.WriteString(runState.UnfinishedDetails.String())
	}
	return result.String()
}
