package undobranches

import (
	"slices"

	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/undo/undodomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	"github.com/git-town/git-town/v22/pkg/set"
)

// BranchChanges describes the changes made to the branches in a Git repo.
// Various types of changes are distinguished.
type BranchChanges struct {
	// Inconsistent changes are changes on both local and tracking branch, but where the local and tracking branch
	// don't have the same SHA before or after.
	// These changes cannot be undone for perennial branches because there is no way to reset the remote branch to the SHA it had before.
	InconsistentlyChanged undodomain.InconsistentChanges
	LocalAdded            gitdomain.LocalBranchNames
	LocalChanged          LocalBranchChange
	LocalRemoved          LocalBranchesSHAs
	LocalRenamed          []LocalBranchRename
	// OmniChanges are changes where the local SHA and the remote SHA are identical before the change as well as after the change,
	OmniChanged LocalBranchChange // a branch had the same SHA locally and remotely, now it has a new SHA locally and remotely, the local and remote SHA are still equal
	// OmniRemoved is when a branch that has the same SHA on its local and tracking branch gets removed.
	OmniRemoved   LocalBranchesSHAs
	RemoteAdded   gitdomain.RemoteBranchNames
	RemoteChanged RemoteBranchChange
	RemoteRemoved RemoteBranchesSHAs
}

// uncomment when debugging the undo logic
// func (self BranchChanges) String() string {
// 	s := strings.Builder{}
// 	s.WriteString("BranchChanges {")
// 	s.WriteString("\n  LocalAdded: ")
// 	s.WriteString(strings.Join(self.LocalAdded.Strings(), ", "))
// 	s.WriteString("\n  LocalRemoved: ")
// 	s.WriteString(strings.Join(self.LocalRemoved.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  LocalChanged: ")
// 	s.WriteString(strings.Join(self.LocalChanged.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  RemoteAdded: ")
// 	s.WriteString(strings.Join(self.RemoteAdded.Strings(), ", "))
// 	s.WriteString("\n  RemoteRemoved: ")
// 	s.WriteString(strings.Join(self.RemoteRemoved.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  RemoteChanged: ")
// 	s.WriteString(strings.Join(self.RemoteChanged.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  OmniRemoved: ")
// 	s.WriteString(strings.Join(self.OmniRemoved.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  OmniChanged: ")
// 	s.WriteString(strings.Join(self.OmniChanged.BranchNames().Strings(), ", "))
// 	s.WriteRune('\n')
// 	return s.String()
// }

// UndoProgram provides the steps to undo the changes described by this BranchChanges instance.
func (self BranchChanges) UndoProgram(args BranchChangesUndoProgramArgs) (undoProgram program.Program, changedBranches gitdomain.LocalBranchNames) {
	result := program.Program{}
	changed := set.Set[gitdomain.LocalBranchName]{}

	omniChanges := CategorizeLocalBranchChange(self.OmniChanged, args.Config)
	changed.Add(omniChanges.Features.BranchNames()...)
	changed.Add(omniChanges.Perennials.BranchNames()...)

	// revert omni-changed perennial branches
	for _, branch := range omniChanges.Perennials.BranchNames() {
		change := omniChanges.Perennials[branch]
		if slices.Contains(args.UndoablePerennialCommits, change.After) {
			result.Add(&opcodes.CheckoutIfNeeded{Branch: branch})
			result.Add(&opcodes.CommitRevertIfNeeded{SHA: change.After})
			if branchInfo, hasBranchInfo := args.BranchInfos.FindByLocalName(branch).Get(); hasBranchInfo {
				if tracking, hasTracking := branchInfo.RemoteName.Get(); hasTracking {
					result.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: branch, TrackingBranch: tracking})
				}
			}
		} else {
			args.FinalMessages.Addf(messages.UndoCannotRevertCommitOnPerennialBranch, change.After)
		}
	}

	// reset omni-changed feature branches
	for _, branch := range omniChanges.Features.BranchNames() {
		change := omniChanges.Features[branch]
		result.Add(&opcodes.CheckoutIfNeeded{Branch: branch})
		result.Add(&opcodes.BranchCurrentResetToSHAIfNeeded{MustHaveSHA: change.After, SetToSHA: change.Before})
		if branchInfo, hasBranchInfo := args.BranchInfos.FindByLocalName(branch).Get(); hasBranchInfo {
			if tracking, hasTracking := branchInfo.RemoteName.Get(); hasTracking {
				result.Add(&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: branch, ForceIfIncludes: true, TrackingBranch: tracking})
			}
		}
	}

	// re-create removed omni-branches
	for _, branch := range self.OmniRemoved.BranchNames() {
		sha := self.OmniRemoved[branch]
		result.Add(&opcodes.BranchCreate{Branch: branch, StartingPoint: sha.Location()})
		result.Add(&opcodes.BranchTrackingCreate{Branch: branch})
	}
	changed.Add(self.OmniRemoved.BranchNames()...)

	inconsistentChanges := CategorizeInconsistentChanges(self.InconsistentlyChanged, args.Config)

	// reset inconsistently changed perennial branches
	for _, inconsistentlyChangedPerennial := range inconsistentChanges.Perennials {
		if omni, isOmni := inconsistentlyChangedPerennial.After.OmniBranch().Get(); isOmni {
			if slices.Contains(args.UndoablePerennialCommits, omni.SHA) {
				result.Add(&opcodes.CheckoutIfNeeded{Branch: omni.Name})
				result.Add(&opcodes.CommitRevertIfNeeded{SHA: omni.SHA})
				if branchInfo, hasBranchInfo := args.BranchInfos.FindByLocalName(omni.Name).Get(); hasBranchInfo {
					if tracking, hasTracking := branchInfo.RemoteName.Get(); hasTracking {
						result.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: omni.Name, TrackingBranch: tracking})
					}
				}
			}
		} else {
			args.FinalMessages.Addf(messages.UndoCannotRevertCommitOnPerennialBranch, inconsistentlyChangedPerennial.After)
		}
		if local, hasLocal := inconsistentlyChangedPerennial.After.Local.Get(); hasLocal {
			changed.Add(local.Name)
		}
		if local, hasLocal := inconsistentlyChangedPerennial.Before.Local.Get(); hasLocal {
			changed.Add(local.Name)
		}
	}

	// reset inconsistently changed feature branches
	for _, inconsistentChange := range inconsistentChanges.Features {
		local, hasLocal := inconsistentChange.Before.Local.Get()
		hasBeforeRemote, beforeRemoteName, beforeRemoteSHA := inconsistentChange.Before.GetRemote()
		AfterSHAs := inconsistentChange.After.GetSHAs()
		if hasLocal && hasBeforeRemote && AfterSHAs.HasBothSHA {
			result.Add(&opcodes.CheckoutIfNeeded{Branch: local.Name})
			result.Add(&opcodes.BranchCurrentResetToSHAIfNeeded{
				MustHaveSHA: AfterSHAs.LocalSHA,
				SetToSHA:    local.SHA,
			})
			result.Add(&opcodes.BranchRemoteSetToSHAIfNeeded{
				Branch:      beforeRemoteName,
				MustHaveSHA: AfterSHAs.RemoteSHA,
				SetToSHA:    beforeRemoteSHA,
			})
		}
		if local, hasLocal := inconsistentChange.After.Local.Get(); hasLocal {
			changed.Add(local.Name)
		}
		if local, hasLocal := inconsistentChange.Before.Local.Get(); hasLocal {
			changed.Add(local.Name)
		}
	}

	// re-create remotely removed feature branches
	removedTrackingBranches := CategorizeRemoteBranchesSHAs(self.RemoteRemoved, args.Config)
	for _, branch := range removedTrackingBranches.Features.BranchNames() {
		sha := removedTrackingBranches.Features[branch]
		result.Add(&opcodes.BranchRemoteCreate{
			Branch: branch.LocalBranchName(),
			SHA:    sha,
		})
	}

	// reset locally changed branches
	for _, localBranch := range self.LocalChanged.BranchNames() {
		change := self.LocalChanged[localBranch]
		result.Add(&opcodes.CheckoutIfNeeded{Branch: localBranch})
		result.Add(&opcodes.BranchCurrentResetToSHAIfNeeded{MustHaveSHA: change.After, SetToSHA: change.Before})
	}
	changed.Add(self.LocalChanged.BranchNames()...)

	// re-create locally removed branches
	for _, removedLocalBranch := range self.LocalRemoved.BranchNames() {
		startingPoint := self.LocalRemoved[removedLocalBranch]
		result.Add(&opcodes.BranchCreate{
			Branch:        removedLocalBranch,
			StartingPoint: startingPoint.Location(),
		})
	}
	changed.Add(self.LocalRemoved.BranchNames()...)

	// restore the name of locally renamed branches
	for _, rename := range self.LocalRenamed {
		result.Add(&opcodes.BranchLocalRename{
			NewName: rename.Before,
			OldName: rename.After,
		})
		changed.Add(rename.After)
		changed.Add(rename.Before)
	}

	// remove locally added branches
	for _, addedLocalBranch := range self.LocalAdded {
		if args.EndBranch == addedLocalBranch {
			result.Add(&opcodes.CheckoutIfNeeded{Branch: args.BeginBranch})
		}
		result.Add(&opcodes.BranchLocalDelete{Branch: addedLocalBranch})
	}
	changed.Add(self.LocalAdded...)

	// Ignore remotely changed perennial branches because we can't force-push to them
	// and we would need the local branch to revert commits on them, but we can't change the local branch.

	// reset remotely changed feature branches
	remoteChanges := CategorizeRemoteBranchChange(self.RemoteChanged, args.Config)
	for _, remoteChangedFeatureBranch := range remoteChanges.Features.BranchNames() {
		change := remoteChanges.Features[remoteChangedFeatureBranch]
		result.Add(&opcodes.BranchRemoteSetToSHAIfNeeded{
			Branch:      remoteChangedFeatureBranch,
			MustHaveSHA: change.After,
			SetToSHA:    change.Before,
		})
	}

	// undo the proposal changes now when the old and new remote branches exist
	result.AddProgram(args.UndoAPIProgram)

	// remove remotely added branches
	for _, addedRemoteBranch := range self.RemoteAdded {
		if addedRemoteBranch.Remote() != gitdomain.RemoteUpstream {
			result.Add(&opcodes.BranchTrackingDelete{
				Branch: addedRemoteBranch,
			})
		}
	}

	// This must be a CheckoutIfExists opcode because this branch might not exist
	// when a Git Town command fails, stores this undo opcode, then gets continued and deletes this branch.
	result.Add(&opcodes.CheckoutIfExists{Branch: args.BeginBranch})
	return result, changed.Values()
}

type BranchChangesUndoProgramArgs struct {
	BeginBranch              gitdomain.LocalBranchName
	BranchInfos              gitdomain.BranchInfos
	Config                   config.ValidatedConfig
	EndBranch                gitdomain.LocalBranchName
	FinalMessages            stringslice.Collector
	UndoAPIProgram           program.Program
	UndoablePerennialCommits []gitdomain.SHA
}
