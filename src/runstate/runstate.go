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
type RunState struct {
	Command                  string                     `json:"Command"`
	IsAbort                  bool                       `exhaustruct:"optional"     json:"IsAbort"`
	IsUndo                   bool                       `exhaustruct:"optional"     json:"IsUndo"`
	AbortSteps               StepList                   `exhaustruct:"optional"     json:"AbortStepList"`
	RunSteps                 StepList                   `json:"RunStepList"`
	UndoSteps                StepList                   `exhaustruct:"optional"     json:"UndoStepList"`
	InitialActiveBranch      domain.LocalBranchName     `json:"InitialActiveBranch"`
	FinalUndoSteps           StepList                   `exhaustruct:"optional"     json:"FinalUndoStepList"`
	UnfinishedDetails        *UnfinishedRunStateDetails `exhaustruct:"optional"     json:"UnfinishedDetails"`
	UndoablePerennialCommits []domain.SHA               `exhaustruct:"optional"     json:"UndoablePerennialCommits"`
}

// AddPushBranchStepAfterCurrentBranchSteps inserts a PushBranchStep
// after all the steps for the current branch.
func (runState *RunState) AddPushBranchStepAfterCurrentBranchSteps(backend *git.BackendCommands) error {
	popped := StepList{}
	for {
		step := runState.RunSteps.Peek()
		if !isCheckoutStep(step) {
			popped.Append(runState.RunSteps.Pop())
		} else {
			currentBranch, err := backend.CurrentBranch()
			if err != nil {
				return err
			}
			runState.RunSteps.Prepend(&steps.PushCurrentBranchStep{CurrentBranch: currentBranch, NoPushHook: false})
			runState.RunSteps.PrependList(popped)
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
	stepList := runState.AbortSteps
	stepList.AppendList(runState.UndoSteps)
	return RunState{
		Command:             runState.Command,
		IsAbort:             true,
		InitialActiveBranch: runState.InitialActiveBranch,
		RunSteps:            stepList,
	}
}

// CreateSkipRunState returns a new Runstate
// that skips operations for the current branch.
func (runState *RunState) CreateSkipRunState() RunState {
	result := RunState{
		Command:             runState.Command,
		InitialActiveBranch: runState.InitialActiveBranch,
		RunSteps:            runState.AbortSteps,
	}
	for _, step := range runState.UndoSteps.List {
		if isCheckoutStep(step) {
			break
		}
		result.RunSteps.Append(step)
	}
	skipping := true
	for _, step := range runState.RunSteps.List {
		if isCheckoutStep(step) {
			skipping = false
		}
		if !skipping {
			result.RunSteps.Append(step)
		}
	}
	result.RunSteps.List = slice.LowerAll[steps.Step](result.RunSteps.List, &steps.RestoreOpenChangesStep{})
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
		RunSteps:                 runState.UndoSteps,
		UndoablePerennialCommits: []domain.SHA{},
	}
	result.RunSteps.Append(&steps.CheckoutStep{Branch: runState.InitialActiveBranch})
	result.RunSteps = result.RunSteps.RemoveDuplicateCheckoutSteps()
	return result
}

func (runState *RunState) HasAbortSteps() bool {
	return !runState.AbortSteps.IsEmpty()
}

func (runState *RunState) HasRunSteps() bool {
	return !runState.RunSteps.IsEmpty()
}

func (runState *RunState) HasUndoSteps() bool {
	return !runState.UndoSteps.IsEmpty()
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
		step := runState.RunSteps.Peek()
		if isCheckoutStep(step) {
			break
		}
		runState.RunSteps.Pop()
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
	result.WriteString(runState.AbortSteps.StringIndented("    "))
	result.WriteString("  RunStepList: ")
	result.WriteString(runState.RunSteps.StringIndented("    "))
	result.WriteString("  UndoStepList: ")
	result.WriteString(runState.UndoSteps.StringIndented("    "))
	if runState.UnfinishedDetails != nil {
		result.WriteString("  UnfineshedDetails: ")
		result.WriteString(runState.UnfinishedDetails.String())
	}
	return result.String()
}
