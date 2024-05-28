package runstate

import (
	"strings"
	"time"

	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// RunState represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has been done so far.
type RunState struct {
	AbortProgram             program.Program `exhaustruct:"optional"`
	BeginBranchesSnapshot    gitdomain.BranchesSnapshot
	BeginConfigSnapshot      undoconfig.ConfigSnapshot
	BeginStashSize           gitdomain.StashSize
	Command                  string
	DryRun                   bool
	EndBranchesSnapshot      Option[gitdomain.BranchesSnapshot]
	EndConfigSnapshot        Option[undoconfig.ConfigSnapshot]
	EndStashSize             Option[gitdomain.StashSize]
	FinalUndoProgram         program.Program `exhaustruct:"optional"`
	RunProgram               program.Program
	UndoablePerennialCommits []gitdomain.SHA                    `exhaustruct:"optional"`
	UnfinishedDetails        OptionP[UnfinishedRunStateDetails] `exhaustruct:"optional"`
}

func EmptyRunState() RunState {
	return RunState{} //exhaustruct:ignore
}

// AddPushBranchAfterCurrentBranchProgram inserts a PushBranch opcode
// after all the opcodes for the current branch.
func (self *RunState) AddPushBranchAfterCurrentBranchProgram(backend git.BackendCommands) error {
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

func (self *RunState) HasAbortProgram() bool {
	return !self.AbortProgram.IsEmpty()
}

func (self *RunState) HasRunProgram() bool {
	return !self.RunProgram.IsEmpty()
}

// IsFinished returns whether or not the run state is unfinished.
func (self *RunState) IsFinished() bool {
	return self.UnfinishedDetails.IsNone()
}

// MarkAsFinished updates the run state to be marked as finished.
func (self *RunState) MarkAsFinished() {
	self.UnfinishedDetails = NoneP[UnfinishedRunStateDetails]()
}

// MarkAsUnfinished updates the run state to be marked as unfinished and populates informational fields.
func (self *RunState) MarkAsUnfinished(backend git.BackendCommands) error {
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
	result.WriteString("\n  AbortProgram: ")
	result.WriteString(self.AbortProgram.StringIndented("    "))
	result.WriteString("  RunProgram: ")
	result.WriteString(self.RunProgram.StringIndented("    "))
	if self.UnfinishedDetails != nil {
		result.WriteString("  UnfineshedDetails: ")
		result.WriteString(self.UnfinishedDetails.String())
	}
	return result.String()
}
