package runstate

import (
	"strings"
	"time"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// RunState represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has been done so far.
type RunState struct {
	AbortProgram             program.Program                            `exhaustruct:"optional"` // opcodes to abort the currently pending Git operation
	BeginBranchesSnapshot    gitdomain.BranchesSnapshot                 // snapshot of the Git branches before the Git Town command that this RunState is for ran
	BeginConfigSnapshot      configdomain.BeginConfigSnapshot           // snapshot of the Git configuration before the Git Town command that this RunState is for ran
	BeginStashSize           gitdomain.StashSize                        // size of the Git stash before the Git Town command that this RunState is for ran
	BranchInfosLastRun       Option[gitdomain.BranchInfos]              // branch infos when the last Git Town command ended
	Command                  string                                     // name of the Git Town command that this RunState is for
	DryRun                   configdomain.DryRun                        // whether the Git Town command that this RunState is for operated in dry-run mode
	EndBranchesSnapshot      Option[gitdomain.BranchesSnapshot]         // snapshot of the Git branches after the Git Town command that this RunState is for ran
	EndConfigSnapshot        Option[configdomain.EndConfigSnapshot]     // snapshot of the Git configuration after the Git Town command that this RunState is for ran
	EndStashSize             Option[gitdomain.StashSize]                // size of the Git stash after the Git Town command that this RunState is for ran
	FinalUndoProgram         program.Program                            `exhaustruct:"optional"` // additional opcodes to run after this RunState was undone
	RunProgram               program.Program                            // remaining opcodes of the Git Town command that this RunState is for
	TouchedBranches          gitdomain.BranchNames                      // the branches that are touched by the Git Town command that this RunState is for
	UndoAPIProgram           program.Program                            // opcodes to undo changes at external systems
	UndoablePerennialCommits []gitdomain.SHA                            `exhaustruct:"optional"` // contains the SHAs of commits on perennial branches that can safely be undone
	UnfinishedDetails        OptionalMutable[UnfinishedRunStateDetails] `exhaustruct:"optional"`
}

func EmptyRunState() RunState {
	return RunState{} //exhaustruct:ignore
}

// inserts a PushBranch opcode after all the opcodes for the current branch
func (self *RunState) AddPushBranchAfterCurrentBranchProgram(gitCommands git.Commands, backend subshelldomain.Querier) error {
	popped := program.Program{}
	for {
		nextOpcode, hasNextOpcode := self.RunProgram.Pop().Get()
		if hasNextOpcode && !opcodes.IsEndOfBranchProgramOpcode(nextOpcode) {
			popped.Add(nextOpcode)
		} else {
			currentBranchOpt, err := gitCommands.CurrentBranch(backend)
			if err != nil {
				return err
			}
			if currentBranch, hasCurrentBranch := currentBranchOpt.Get(); hasCurrentBranch {
				self.RunProgram.Prepend(&opcodes.PushCurrentBranchIfNeeded{
					CurrentBranch: currentBranch,
				})
			}
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
func (self *RunState) MarkAsFinished(endBranchesSnapshot gitdomain.BranchesSnapshot) {
	self.UnfinishedDetails = MutableNone[UnfinishedRunStateDetails]()
	self.EndBranchesSnapshot = Some(endBranchesSnapshot)
}

// MarkAsUnfinished updates the run state to be marked as unfinished and populates informational fields.
func (self *RunState) MarkAsUnfinished(gitCommands git.Commands, backend subshelldomain.Querier, canSkip bool) error {
	currentBranchOpt, err := gitCommands.CurrentBranch(backend)
	if err != nil {
		return err
	}
	if currentBranch, hasCurrentBranch := currentBranchOpt.Get(); hasCurrentBranch {
		self.UnfinishedDetails = MutableSome(&UnfinishedRunStateDetails{
			CanSkip:   canSkip,
			EndBranch: currentBranch,
			EndTime:   time.Now(),
		})
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
		opcode, hasOpcode := self.RunProgram.Pop().Get()
		if !hasOpcode || opcodes.IsEndOfBranchProgramOpcode(opcode) {
			break
		}
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
		result.WriteString(messages.DialogResultNone)
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
		result.WriteString(branchInfo.LocalSHA.StringOr(messages.DialogResultNone))
		result.WriteString("\n      Remote: ")
		result.WriteString(branchInfo.RemoteSHA.StringOr(messages.DialogResultNone))
		result.WriteRune('\n')
	}
}
