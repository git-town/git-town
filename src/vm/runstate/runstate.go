package runstate

import (
	"strconv"
	"strings"
	"time"

	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// RunState represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has been done so far.
type RunState struct {
	AbortProgram             program.Program `exhaustruct:"optional"`
	Command                  string
	DryRun                   bool
	FinalUndoProgram         program.Program `exhaustruct:"optional"`
	InitialActiveBranch      gitdomain.LocalBranchName
	IsUndo                   bool `exhaustruct:"optional"`
	RunProgram               program.Program
	UndoProgram              program.Program            `exhaustruct:"optional"`
	UndoablePerennialCommits []gitdomain.SHA            `exhaustruct:"optional"`
	UnfinishedDetails        *UnfinishedRunStateDetails `exhaustruct:"optional"`
}

func EmptyRunState() RunState {
	return RunState{} //nolint:exhaustruct
}

// AddPushBranchAfterCurrentBranchProgram inserts a PushBranch opcode
// after all the opcodes for the current branch.
func (self *RunState) AddPushBranchAfterCurrentBranchProgram(backend *git.BackendCommands) error {
	popped := program.Program{}
	for {
		nextOpcode := self.RunProgram.Peek()
		if !shared.IsEndOfBranchProgramOpcode(nextOpcode) {
			popped.Add(self.RunProgram.Pop())
		} else {
			currentBranch, err := backend.CurrentBranch()
			if err != nil {
				return err
			}
			self.RunProgram.Prepend(&opcodes.PushCurrentBranch{CurrentBranch: currentBranch})
			self.RunProgram.PrependProgram(popped)
			break
		}
	}
	return nil
}

// CreateAbortRunState returns a new runstate
// to be run to aborting and undoing the Git Town command
// represented by this runstate.
func (self *RunState) CreateAbortRunState() RunState {
	abortProgram := self.AbortProgram
	abortProgram.AddProgram(self.UndoProgram)
	return RunState{
		Command:             self.Command,
		DryRun:              self.DryRun,
		InitialActiveBranch: self.InitialActiveBranch,
		IsUndo:              true,
		RunProgram:          abortProgram,
	}
}

// CreateSkipRunState returns a new Runstate
// that skips operations for the current branch.
func (self *RunState) CreateSkipRunState() RunState {
	result := RunState{
		Command:             self.Command,
		DryRun:              self.DryRun,
		InitialActiveBranch: self.InitialActiveBranch,
		RunProgram:          self.AbortProgram,
	}
	// undo the operations done on the current branch so far
	// by copying the respective undo-opcodes into the runprogram
	for _, opcode := range self.UndoProgram {
		if shared.IsCheckoutOpcode(opcode) {
			break
		}
		result.RunProgram.Add(opcode)
	}
	// skip the remaining run-opcodes for this branch
	skipping := true
	for _, opcode := range self.RunProgram {
		if shared.IsEndOfBranchProgramOpcode(opcode) {
			skipping = false
		}
		if !skipping {
			result.RunProgram.Add(opcode)
		}
	}
	result.RunProgram.MoveToEnd(&opcodes.RestoreOpenChanges{})
	return result
}

// CreateUndoRunState returns a new runstate
// to be run when undoing the Git Town command
// represented by this runstate.
func (self *RunState) CreateUndoRunState() RunState {
	result := RunState{
		Command:                  self.Command,
		DryRun:                   self.DryRun,
		InitialActiveBranch:      self.InitialActiveBranch,
		IsUndo:                   true,
		RunProgram:               self.UndoProgram,
		UndoablePerennialCommits: []gitdomain.SHA{},
	}
	result.RunProgram.Add(&opcodes.Checkout{Branch: self.InitialActiveBranch})
	result.RunProgram.RemoveDuplicateCheckout()
	return result
}

func (self *RunState) HasAbortProgram() bool {
	return !self.AbortProgram.IsEmpty()
}

func (self *RunState) HasRunProgram() bool {
	return !self.RunProgram.IsEmpty()
}

func (self *RunState) HasUndoProgram() bool {
	return !self.UndoProgram.IsEmpty()
}

// IsFinished returns whether or not the run state is unfinished.
func (self *RunState) IsFinished() bool {
	return self.UnfinishedDetails == nil
}

// MarkAsFinished updates the run state to be marked as finished.
func (self *RunState) MarkAsFinished() {
	self.UnfinishedDetails = nil
}

// MarkAsUnfinished updates the run state to be marked as unfinished and populates informational fields.
func (self *RunState) MarkAsUnfinished(backend *git.BackendCommands) error {
	currentBranch, err := backend.CurrentBranch()
	if err != nil {
		return err
	}
	self.UnfinishedDetails = &UnfinishedRunStateDetails{
		CanSkip:   false,
		EndBranch: currentBranch,
		EndTime:   time.Now(),
	}
	return nil
}

// RegisterUndoablePerennialCommit stores the given commit on a perennial branch as undoable.
// This method is used as a callback.
func (self *RunState) RegisterUndoablePerennialCommit(commit gitdomain.SHA) {
	self.UndoablePerennialCommits = append(self.UndoablePerennialCommits, commit)
}

// SkipCurrentBranchProgram removes the opcodes for the current branch
// from this run state.
func (self *RunState) SkipCurrentBranchProgram() {
	for {
		opcode := self.RunProgram.Peek()
		if shared.IsEndOfBranchProgramOpcode(opcode) {
			break
		}
		self.RunProgram.Pop()
	}
}

func (self *RunState) String() string {
	result := strings.Builder{}
	result.WriteString("RunState:\n")
	result.WriteString("  Command: ")
	result.WriteString(self.Command)
	result.WriteString("\n  IsUndo: ")
	result.WriteString(strconv.FormatBool(self.IsUndo))
	result.WriteString("\n  AbortProgram: ")
	result.WriteString(self.AbortProgram.StringIndented("    "))
	result.WriteString("  RunProgram: ")
	result.WriteString(self.RunProgram.StringIndented("    "))
	result.WriteString("  UndoProgram: ")
	result.WriteString(self.UndoProgram.StringIndented("    "))
	if self.UnfinishedDetails != nil {
		result.WriteString("  UnfineshedDetails: ")
		result.WriteString(self.UnfinishedDetails.String())
	}
	return result.String()
}
