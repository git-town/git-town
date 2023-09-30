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
func (rs *RunState) AddPushBranchStepAfterCurrentBranchSteps(backend *git.BackendCommands) error {
	popped := StepList{}
	for {
		step := rs.RunStepList.Peek()
		if !isCheckoutStep(step) {
			popped.Append(rs.RunStepList.Pop())
		} else {
			currentBranch, err := backend.CurrentBranch()
			if err != nil {
				return err
			}
			rs.RunStepList.Prepend(&steps.PushCurrentBranchStep{CurrentBranch: currentBranch, NoPushHook: false})
			rs.RunStepList.PrependList(popped)
			break
		}
	}
	return nil
}

// RegisterUndoablePerennialCommit stores the given commit on a perennial branch as undoable.
// This method is used as a callback.
func (rs *RunState) RegisterUndoablePerennialCommit(commit domain.SHA) {
	rs.UndoablePerennialCommits = append(rs.UndoablePerennialCommits, commit)
}

// CreateAbortRunState returns a new runstate
// to be run to aborting and undoing the Git Town command
// represented by this runstate.
func (rs *RunState) CreateAbortRunState() RunState {
	stepList := rs.AbortStepList
	stepList.AppendList(rs.UndoStepList)
	return RunState{
		Command:             rs.Command,
		IsAbort:             true,
		InitialActiveBranch: rs.InitialActiveBranch,
		RunStepList:         stepList,
	}
}

// CreateSkipRunState returns a new Runstate
// that skips operations for the current branch.
func (rs *RunState) CreateSkipRunState() RunState {
	result := RunState{
		Command:             rs.Command,
		InitialActiveBranch: rs.InitialActiveBranch,
		RunStepList:         rs.AbortStepList,
	}
	for _, step := range rs.UndoStepList.List {
		if isCheckoutStep(step) {
			break
		}
		result.RunStepList.Append(step)
	}
	skipping := true
	for _, step := range rs.RunStepList.List {
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
func (rs *RunState) CreateUndoRunState() RunState {
	result := RunState{
		Command:                  rs.Command,
		InitialActiveBranch:      rs.InitialActiveBranch,
		IsUndo:                   true,
		RunStepList:              rs.UndoStepList,
		UndoablePerennialCommits: []domain.SHA{},
	}
	result.RunStepList.Append(&steps.CheckoutStep{Branch: rs.InitialActiveBranch})
	result.RunStepList = result.RunStepList.RemoveDuplicateCheckoutSteps()
	return result
}

func (rs *RunState) HasAbortSteps() bool {
	return !rs.AbortStepList.IsEmpty()
}

func (rs *RunState) HasRunSteps() bool {
	return !rs.RunStepList.IsEmpty()
}

func (rs *RunState) HasUndoSteps() bool {
	return !rs.UndoStepList.IsEmpty()
}

// IsUnfinished returns whether or not the run state is unfinished.
func (rs *RunState) IsUnfinished() bool {
	return rs.UnfinishedDetails != nil
}

// MarkAsFinished updates the run state to be marked as finished.
func (rs *RunState) MarkAsFinished() {
	rs.UnfinishedDetails = nil
}

// MarkAsUnfinished updates the run state to be marked as unfinished and populates informational fields.
func (rs *RunState) MarkAsUnfinished(backend *git.BackendCommands) error {
	currentBranch, err := backend.CurrentBranch()
	if err != nil {
		return err
	}
	rs.UnfinishedDetails = &UnfinishedRunStateDetails{
		CanSkip:   false,
		EndBranch: currentBranch,
		EndTime:   time.Now(),
	}
	return nil
}

// SkipCurrentBranchSteps removes the steps for the current branch
// from this run state.
func (rs *RunState) SkipCurrentBranchSteps() {
	for {
		step := rs.RunStepList.Peek()
		if isCheckoutStep(step) {
			break
		}
		rs.RunStepList.Pop()
	}
}

func (rs *RunState) String() string {
	result := strings.Builder{}
	result.WriteString("RunState:\n")
	result.WriteString("  Command: ")
	result.WriteString(rs.Command)
	result.WriteString("\n  IsAbort: ")
	result.WriteString(fmt.Sprintf("%t", rs.IsAbort))
	result.WriteString("\n  IsUndo: ")
	result.WriteString(fmt.Sprintf("%t", rs.IsUndo))
	result.WriteString("\n  AbortStepList: ")
	result.WriteString(rs.AbortStepList.StringIndented("    "))
	result.WriteString("  RunStepList: ")
	result.WriteString(rs.RunStepList.StringIndented("    "))
	result.WriteString("  UndoStepList: ")
	result.WriteString(rs.UndoStepList.StringIndented("    "))
	if rs.UnfinishedDetails != nil {
		result.WriteString("  UnfineshedDetails: ")
		result.WriteString(rs.UnfinishedDetails.String())
	}
	return result.String()
}
