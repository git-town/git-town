package vmstate

import (
	"strings"
	"time"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Data represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has been done so far.
type Data struct {
	AbortProgram             program.Program                    `exhaustruct:"optional"` // opcodes to abort the currently pending Git operation
	BeginBranchesSnapshot    gitdomain.BranchesSnapshot         // snapshot of the Git branches before the Git Town command that this RunState is for ran
	BeginConfigSnapshot      undoconfig.ConfigSnapshot          // snapshot of the Git configuration before the Git Town command that this RunState is for ran
	BeginStashSize           gitdomain.StashSize                // size of the Git stash before the Git Town command that this RunState is for ran
	BranchInfosLastRun       Option[gitdomain.BranchInfos]      // branch infos when the last Git Town command ended
	Command                  string                             // name of the Git Town command that this RunState is for
	DryRun                   configdomain.DryRun                // whether the Git Town command that this RunState is for operated in dry-run mode
	EndBranchesSnapshot      Option[gitdomain.BranchesSnapshot] // snapshot of the Git branches after the Git Town command that this RunState is for ran
	EndConfigSnapshot        Option[undoconfig.ConfigSnapshot]  // snapshot of the Git configuration after the Git Town command that this RunState is for ran
	EndStashSize             Option[gitdomain.StashSize]        // size of the Git stash after the Git Town command that this RunState is for ran
	FinalUndoProgram         program.Program                    `exhaustruct:"optional"` // additional opcodes to run after this RunState was undone
	RunProgram               program.Program                    // remaining opcodes of the Git Town command that this RunState is for
	TouchedBranches          gitdomain.BranchNames              // the branches that are touched by the Git Town command that this RunState is for
	UndoAPIProgram           program.Program                    // opcodes to undo changes at external systems
	UndoablePerennialCommits []gitdomain.SHA                    `exhaustruct:"optional"` // contains the SHAs of commits on perennial branches that can safely be undone
	UnfinishedDetails        OptionalMutable[UnfinishedData]    `exhaustruct:"optional"`
}

func EmptyRunState() Data {
	return Data{} //exhaustruct:ignore
}

// inserts a PushBranch opcode after all the opcodes for the current branch
func (self *Data) AddPushBranchAfterCurrentBranchProgram(gitCommands git.Commands, backend gitdomain.Querier) error {
	popped := program.Program{}
	for {
		nextOpcode := self.RunProgram.Peek()
		if !opcodes.IsEndOfBranchProgramOpcode(nextOpcode) {
			popped.Add(self.RunProgram.Pop())
		} else {
			currentBranch, err := gitCommands.CurrentBranch(backend)
			if err != nil {
				return err
			}
			self.RunProgram.Prepend(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: currentBranch})
			self.RunProgram.PrependProgram(popped)
			break
		}
	}
	return nil
}

func (self *Data) HasAbortProgram() bool {
	return !self.AbortProgram.IsEmpty()
}

func (self *Data) HasRunProgram() bool {
	return !self.RunProgram.IsEmpty()
}

// IsFinished returns whether or not the run state is unfinished.
func (self *Data) IsFinished() bool {
	return self.UnfinishedDetails.IsNone()
}

// MarkAsFinished updates the run state to be marked as finished.
func (self *Data) MarkAsFinished(endBranchesSnapshot gitdomain.BranchesSnapshot) {
	self.UnfinishedDetails = MutableNone[UnfinishedData]()
	self.EndBranchesSnapshot = Some(endBranchesSnapshot)
}

// MarkAsUnfinished updates the run state to be marked as unfinished and populates informational fields.
func (self *Data) MarkAsUnfinished(gitCommands git.Commands, backend gitdomain.Querier, canSkip bool) error {
	currentBranch, err := gitCommands.CurrentBranch(backend)
	if err != nil {
		return err
	}
	self.UnfinishedDetails = MutableSome(&UnfinishedData{
		CanSkip:   canSkip,
		EndBranch: currentBranch,
		EndTime:   time.Now(),
	})
	return nil
}

// RegisterUndoablePerennialCommit stores the given commit on a perennial branch as undoable.
// This method is used as a callback.
func (self *Data) RegisterUndoablePerennialCommit(commit gitdomain.SHA) {
	self.UndoablePerennialCommits = append(self.UndoablePerennialCommits, commit)
}

// SkipCurrentBranchProgram removes the opcodes for the current branch
// from this run state.
func (self *Data) SkipCurrentBranchProgram() {
	for {
		opcode := self.RunProgram.Peek()
		if opcodes.IsEndOfBranchProgramOpcode(opcode) {
			break
		}
		self.RunProgram.Pop()
	}
}

func (self *Data) String() string {
	result := strings.Builder{}
	result.WriteString("RunState:\n")
	result.WriteString("  Command: ")
	result.WriteString(self.Command)
	result.WriteString("\n  AbortProgram: ")
	result.WriteString(self.AbortProgram.StringIndented("    "))
	result.WriteString("  RunProgram: ")
	result.WriteString(self.RunProgram.StringIndented("    "))
	if unfinishedDetails, hasUnfinishedDetails := self.UnfinishedDetails.Get(); hasUnfinishedDetails {
		result.WriteString("  UnfinishedDetails: ")
		result.WriteString(unfinishedDetails.String())
	}
	result.WriteString("  Touched branches: ")
	result.WriteString(self.TouchedBranches.Join(", "))
	result.WriteString("\n  Before snapshot: \n")
	writeBranchInfos(&result, self.BeginBranchesSnapshot.Branches)
	result.WriteString("\n  After snapshot: \n")
	if endSnapshot, has := self.EndBranchesSnapshot.Get(); has {
		writeBranchInfos(&result, endSnapshot.Branches)
	} else {
		result.WriteString("(none)")
	}
	return result.String()
}

func writeBranchInfos(result *strings.Builder, branchInfos gitdomain.BranchInfos) {
	for _, branchInfo := range branchInfos {
		result.WriteString("    Branch: ")
		result.WriteString(branchInfo.GetLocalOrRemoteName().String())
		result.WriteString(" (")
		result.WriteString(string(branchInfo.SyncStatus))
		result.WriteString(")\n      Local: ")
		result.WriteString(branchInfo.LocalSHA.StringOr("(none)"))
		result.WriteString("\n      Remote: ")
		result.WriteString(branchInfo.RemoteSHA.StringOr("(none)"))
		result.WriteRune('\n')
	}
}
