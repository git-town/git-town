package runstate

import (
	"strings"
	"time"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/undo/undoconfig"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// RunState represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has been done so far.
type RunState struct {
	AbortProgram             program.Program                    `exhaustruct:"optional"` // opcodes to abort the currently pending Git operation
	BeginBranchesSnapshot    gitdomain.BranchesSnapshot         // snapshot of the Git branches before the Git Town command that this RunState is for ran
	BeginConfigSnapshot      undoconfig.ConfigSnapshot          // snapshot of the Git configuration before the Git Town command that this RunState is for ran
	BeginStashSize           gitdomain.StashSize                // size of the Git stash before the Git Town command that this RunState is for ran
	Command                  string                             // name of the Git Town command that this RunState is for
	DryRun                   configdomain.DryRun                // whether the Git Town command that this RunState is for operated in dry-run mode
	EndBranchesSnapshot      Option[gitdomain.BranchesSnapshot] // snapshot of the Git branches after the Git Town command that this RunState is for ran
	EndConfigSnapshot        Option[undoconfig.ConfigSnapshot]  // snapshot of the Git configuration after the Git Town command that this RunState is for ran
	EndStashSize             Option[gitdomain.StashSize]        // size of the Git stash after the Git Town command that this RunState is for ran
	FinalUndoProgram         program.Program                    `exhaustruct:"optional"` // additional opcodes to run after this RunState was undone
	RunProgram               program.Program                    // remaining opcodes of the Git Town command that this RunState is for
	TouchedBranches          []gitdomain.BranchName             // the branches that are touched by the Git Town command that this RunState is for
	UndoablePerennialCommits []gitdomain.SHA                    `exhaustruct:"optional"` // contains the SHAs of commits on perennial branches that can safely be undone
	UnfinishedDetails        OptionP[UnfinishedRunStateDetails] `exhaustruct:"optional"`
}

func EmptyRunState() RunState {
	return RunState{} //exhaustruct:ignore
}

// inserts a PushBranch opcode after all the opcodes for the current branch
func (self *RunState) AddPushBranchAfterCurrentBranchProgram(gitCommands git.Commands, backend gitdomain.Querier) error {
	popped := program.Program{}
	for {
		nextOpcode := self.RunProgram.Peek()
		if !shared.IsEndOfBranchProgramOpcode(nextOpcode) {
			popped.Add(self.RunProgram.Pop())
		} else {
			currentBranch, err := gitCommands.CurrentBranch(backend)
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
func (self *RunState) MarkAsUnfinished(gitCommands git.Commands, backend gitdomain.Querier) error {
	currentBranch, err := gitCommands.CurrentBranch(backend)
	if err != nil {
		return err
	}
	self.UnfinishedDetails = SomeP(&UnfinishedRunStateDetails{
		CanSkip:   false,
		EndBranch: currentBranch,
		EndTime:   time.Now(),
	})
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
	if unfinishedDetails, hasUnfinishedDetails := self.UnfinishedDetails.Get(); hasUnfinishedDetails {
		result.WriteString("  UnfineshedDetails: ")
		result.WriteString(unfinishedDetails.String())
	}
	return result.String()
}
