package runstate

import (
	"fmt"
	"strings"
	"time"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/program"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// RunState represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has been done so far.
type RunState struct {
	Command                  string                     `json:"Command"`
	IsAbort                  bool                       `exhaustruct:"optional"     json:"IsAbort"`
	IsUndo                   bool                       `exhaustruct:"optional"     json:"IsUndo"`
	AbortProgram             program.Program            `exhaustruct:"optional"     json:"AbortProgram"`
	RunProgram               program.Program            `json:"RunProgram"`
	UndoProgram              program.Program            `exhaustruct:"optional"     json:"UndoProgram"`
	InitialActiveBranch      domain.LocalBranchName     `json:"InitialActiveBranch"`
	FinalUndoProgram         program.Program            `exhaustruct:"optional"     json:"FinalUndoProgram"`
	UnfinishedDetails        *UnfinishedRunStateDetails `exhaustruct:"optional"     json:"UnfinishedDetails"`
	UndoablePerennialCommits []domain.SHA               `exhaustruct:"optional"     json:"UndoablePerennialCommits"`
}

// AddPushBranchStepAfterCurrentBranchProgram inserts a PushBranchStep
// after all the steps for the current branch.
func (rs *RunState) AddPushBranchStepAfterCurrentBranchProgram(backend *git.BackendCommands) error {
	popped := program.Program{}
	for {
		nextStep := rs.RunProgram.Peek()
		if !program.IsCheckoutStep(nextStep) {
			popped.Add(rs.RunProgram.Pop())
		} else {
			currentBranch, err := backend.CurrentBranch()
			if err != nil {
				return err
			}
			rs.RunProgram.Prepend(&opcode.PushCurrentBranch{CurrentBranch: currentBranch, NoPushHook: false})
			rs.RunProgram.PrependProgram(popped)
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
	abortProgram := rs.AbortProgram
	abortProgram.AddProgram(rs.UndoProgram)
	return RunState{
		Command:             rs.Command,
		IsAbort:             true,
		InitialActiveBranch: rs.InitialActiveBranch,
		RunProgram:          abortProgram,
	}
}

// CreateSkipRunState returns a new Runstate
// that skips operations for the current branch.
func (rs *RunState) CreateSkipRunState() RunState {
	result := RunState{
		Command:             rs.Command,
		InitialActiveBranch: rs.InitialActiveBranch,
		RunProgram:          rs.AbortProgram,
	}
	for _, step := range rs.UndoProgram.Opcodes {
		if program.IsCheckoutStep(step) {
			break
		}
		result.RunProgram.Add(step)
	}
	skipping := true
	for _, step := range rs.RunProgram.Opcodes {
		if program.IsCheckoutStep(step) {
			skipping = false
		}
		if !skipping {
			result.RunProgram.Add(step)
		}
	}
	result.RunProgram.Opcodes = slice.LowerAll[shared.Opcode](result.RunProgram.Opcodes, &opcode.RestoreOpenChanges{})
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
		RunProgram:               rs.UndoProgram,
		UndoablePerennialCommits: []domain.SHA{},
	}
	result.RunProgram.Add(&opcode.Checkout{Branch: rs.InitialActiveBranch})
	result.RunProgram = result.RunProgram.RemoveDuplicateCheckoutSteps()
	return result
}

func (rs *RunState) HasAbortProgram() bool {
	return !rs.AbortProgram.IsEmpty()
}

func (rs *RunState) HasRunProgram() bool {
	return !rs.RunProgram.IsEmpty()
}

func (rs *RunState) HasUndoProgram() bool {
	return !rs.UndoProgram.IsEmpty()
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

// SkipCurrentBranchProgram removes the steps for the current branch
// from this run state.
func (rs *RunState) SkipCurrentBranchProgram() {
	for {
		step := rs.RunProgram.Peek()
		if program.IsCheckoutStep(step) {
			break
		}
		rs.RunProgram.Pop()
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
	result.WriteString("\n  AbortProgram: ")
	result.WriteString(rs.AbortProgram.StringIndented("    "))
	result.WriteString("  RunProgram: ")
	result.WriteString(rs.RunProgram.StringIndented("    "))
	result.WriteString("  UndoProgram: ")
	result.WriteString(rs.UndoProgram.StringIndented("    "))
	if rs.UnfinishedDetails != nil {
		result.WriteString("  UnfineshedDetails: ")
		result.WriteString(rs.UnfinishedDetails.String())
	}
	return result.String()
}
