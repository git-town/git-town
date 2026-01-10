package undobranches

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	"github.com/git-town/git-town/v22/internal/undo/undodomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// BranchSpans describes how a Git Town command has modified the branches in a Git repository.
type BranchSpans []BranchSpan

func NewBranchSpans(beforeSnapshot, afterSnapshot gitdomain.BranchesSnapshot) BranchSpans {
	result := BranchSpans{}
	for _, before := range beforeSnapshot.Branches {
		after := afterSnapshot.Branches.FindMatchingRecord(before)
		result = append(result, BranchSpan{
			Before: Some(before),
			After:  after.ToOption(),
		})
	}
	for _, after := range afterSnapshot.Branches {
		if beforeSnapshot.Branches.FindMatchingRecord(after).IsNone() {
			result = append(result, BranchSpan{
				Before: None[gitdomain.BranchInfo](),
				After:  Some(after),
			})
		}
	}
	return result
}

// Changes describes the specific changes made in this BranchSpans.
func (self BranchSpans) Changes() BranchChanges {
	inconsistentlyChanged := undodomain.InconsistentChanges{}
	localAdded := gitdomain.LocalBranchNames{}
	localChanged := LocalBranchChange{}
	localRemoved := LocalBranchesSHAs{}
	localRenamed := []LocalBranchRename{}
	omniChanged := LocalBranchChange{}
	omniRemoved := LocalBranchesSHAs{}
	remoteAdded := gitdomain.RemoteBranchNames{}
	remoteChanged := map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{}
	remoteRemoved := map[gitdomain.RemoteBranchName]gitdomain.SHA{}
	for _, branchSpan := range self {
		if omniRemoveData := branchSpan.IsOmniRemove(); omniRemoveData.IsOmniRemove {
			omniRemoved[omniRemoveData.Name] = omniRemoveData.SHA
			continue
		}
		if omniChangeData := branchSpan.IsOmniChange(); omniChangeData.IsOmniChange {
			omniChanged[omniChangeData.Name] = undodomain.Change[gitdomain.SHA]{
				Before: omniChangeData.SHABefore,
				After:  omniChangeData.SHAAfter,
			}
			continue
		}
		isLocalRename, beforeName, afterName := branchSpan.IsLocalRename()
		if isLocalRename {
			localRenamed = append(localRenamed, LocalBranchRename{
				After:  afterName,
				Before: beforeName,
			})
		}
		isInconsistentChange, before, after := branchSpan.IsInconsistentChange()
		if isInconsistentChange {
			inconsistentlyChanged = append(inconsistentlyChanged, undodomain.InconsistentChange{
				Before: before,
				After:  after,
			})
			continue
		}
		if localAddedData := branchSpan.LocalAdded(); localAddedData.IsAdded {
			localAdded = append(localAdded, localAddedData.Name)
		} else if localRemovedResult := branchSpan.LocalRemoved(); localRemovedResult.IsRemoved {
			localRemoved[localRemovedResult.Name] = localRemovedResult.SHA
		} else if localChangedResult := branchSpan.LocalChanged(); localChangedResult.IsChanged {
			localChanged[localChangedResult.Name] = undodomain.Change[gitdomain.SHA]{
				Before: localChangedResult.SHABefore,
				After:  localChangedResult.SHAAfter,
			}
		}
		if remoteAddedResult := branchSpan.RemoteAdded(); remoteAddedResult.IsAdded {
			remoteAdded = append(remoteAdded, remoteAddedResult.Name)
		} else if remoteRemovedResult := branchSpan.RemoteRemoved(); remoteRemovedResult.IsRemoved {
			remoteRemoved[remoteRemovedResult.Name] = remoteRemovedResult.SHA
		} else if remoteChangedResult := branchSpan.RemoteChanged(); remoteChangedResult.IsChanged {
			remoteChanged[remoteAddedResult.Name] = undodomain.Change[gitdomain.SHA]{
				Before: remoteChangedResult.SHABefore,
				After:  remoteChangedResult.SHAAfter,
			}
		}
	}
	return BranchChanges{
		InconsistentlyChanged: inconsistentlyChanged,
		LocalAdded:            localAdded,
		LocalChanged:          localChanged,
		LocalRemoved:          localRemoved,
		LocalRenamed:          localRenamed,
		OmniChanged:           omniChanged,
		OmniRemoved:           omniRemoved,
		RemoteAdded:           remoteAdded,
		RemoteChanged:         remoteChanged,
		RemoteRemoved:         remoteRemoved,
	}
}

// keeps only the branch spans that contain any of the given branches
func (self BranchSpans) KeepOnly(branchesToKeep []gitdomain.BranchName) BranchSpans {
	result := BranchSpans{}
	for _, branchSpan := range self {
		if slice.ContainsAny(branchSpan.BranchNames(), branchesToKeep) {
			result = append(result, branchSpan)
		}
	}
	return result
}
