package state

import (
	"fmt"
	"strings"
	"time"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/steps"
)

// RunState represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has been done so far.
type RunState struct {
	Command                  string                     `json:"Command"`
	IsAbort                  bool                       `exhaustruct:"optional"     json:"IsAbort"`
	IsUndo                   bool                       `exhaustruct:"optional"     json:"IsUndo"`
	AbortSteps               steps.List                 `exhaustruct:"optional"     json:"AbortSteps"`
	RunSteps                 steps.List                 `json:"RunSteps"`
	UndoSteps                steps.List                 `exhaustruct:"optional"     json:"UndoSteps"`
	InitialActiveBranch      domain.LocalBranchName     `json:"InitialActiveBranch"`
	FinalUndoSteps           steps.List                 `exhaustruct:"optional"     json:"FinalUndoSteps"`
	UnfinishedDetails        *UnfinishedRunStateDetails `exhaustruct:"optional"     json:"UnfinishedDetails"`
	UndoablePerennialCommits []domain.SHA               `exhaustruct:"optional"     json:"UndoablePerennialCommits"`
}

// AddPushBranchStepAfterCurrentBranchSteps inserts a PushBranchStep
// after all the steps for the current branch.
func (rs *RunState) AddPushBranchStepAfterCurrentBranchSteps(backend *git.BackendCommands) error {
	popped := steps.List{}
	for {
		nextStep := rs.RunSteps.Peek()
		if !steps.IsCheckoutStep(nextStep) {
			popped.Add(rs.RunSteps.Pop())
		} else {
			currentBranch, err := backend.CurrentBranch()
			if err != nil {
				return err
			}
			rs.RunSteps.Prepend(&step.PushCurrentBranch{CurrentBranch: currentBranch, NoPushHook: false})
			rs.RunSteps.PrependList(popped)
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
	stepList := rs.AbortSteps
	stepList.AddList(rs.UndoSteps)
	return RunState{
		Command:             rs.Command,
		IsAbort:             true,
		InitialActiveBranch: rs.InitialActiveBranch,
		RunSteps:            stepList,
	}
}

// CreateSkipRunState returns a new Runstate
// that skips operations for the current branch.
func (rs *RunState) CreateSkipRunState() RunState {
	result := RunState{
		Command:             rs.Command,
		InitialActiveBranch: rs.InitialActiveBranch,
		RunSteps:            rs.AbortSteps,
	}
	for _, step := range rs.UndoSteps.List {
		if steps.IsCheckoutStep(step) {
			break
		}
		result.RunSteps.Add(step)
	}
	skipping := true
	for _, step := range rs.RunSteps.List {
		if steps.IsCheckoutStep(step) {
			skipping = false
		}
		if !skipping {
			result.RunSteps.Add(step)
		}
	}
	result.RunSteps.List = slice.LowerAll[step.Step](result.RunSteps.List, &step.RestoreOpenChanges{})
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
		RunSteps:                 rs.UndoSteps,
		UndoablePerennialCommits: []domain.SHA{},
	}
	result.RunSteps.Add(&step.Checkout{Branch: rs.InitialActiveBranch})
	result.RunSteps = result.RunSteps.RemoveDuplicateCheckoutSteps()
	return result
}

func (rs *RunState) HasAbortSteps() bool {
	return !rs.AbortSteps.IsEmpty()
}

func (rs *RunState) HasRunSteps() bool {
	return !rs.RunSteps.IsEmpty()
}

func (rs *RunState) HasUndoSteps() bool {
	return !rs.UndoSteps.IsEmpty()
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
		step := rs.RunSteps.Peek()
		if steps.IsCheckoutStep(step) {
			break
		}
		rs.RunSteps.Pop()
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
	result.WriteString(rs.AbortSteps.StringIndented("    "))
	result.WriteString("  RunStepList: ")
	result.WriteString(rs.RunSteps.StringIndented("    "))
	result.WriteString("  UndoStepList: ")
	result.WriteString(rs.UndoSteps.StringIndented("    "))
	if rs.UnfinishedDetails != nil {
		result.WriteString("  UnfineshedDetails: ")
		result.WriteString(rs.UnfinishedDetails.String())
	}
	return result.String()
}
