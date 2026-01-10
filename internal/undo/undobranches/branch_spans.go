package undobranches

import (
	"maps"

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
		if omniRemove, isOmniRemove := branchSpan.IsOmniRemove().Get(); isOmniRemove {
			maps.Copy(omniRemoved, omniRemove)
			continue
		}
		if omniChange, hasOmniChange := branchSpan.OmniChange().Get(); hasOmniChange {
			maps.Copy(omniChanged, omniChange)
			continue
		}
		localBranchRename, isLocalBranchRename := branchSpan.LocalRename().Get()
		if isLocalBranchRename {
			localRenamed = append(localRenamed, localBranchRename)
		}
		inconsistentChange, isInconsistentChange := branchSpan.InconsistentChange().Get()
		if isInconsistentChange {
			inconsistentlyChanged = append(inconsistentlyChanged, undodomain.InconsistentChange{
				Before: inconsistentChange.Before,
				After:  inconsistentChange.After,
			})
			continue
		}
		if localAddedBranch, isLocalAdded := branchSpan.LocalAdded().Get(); isLocalAdded {
			localAdded = append(localAdded, localAddedBranch)
		} else if isLocalRemoved, beforeBranch, beforeSHA := branchSpan.LocalRemoved(); isLocalRemoved {
			localRemoved[beforeBranch] = beforeSHA
		} else if isLocalChanged, branch, beforeSHA, afterSHA := branchSpan.LocalChanged(); isLocalChanged {
			localChanged[branch] = undodomain.Change[gitdomain.SHA]{
				Before: beforeSHA,
				After:  afterSHA,
			}
		}
		if isRemoteAdded, remoteBranchName, _ := branchSpan.RemoteAdded(); isRemoteAdded {
			remoteAdded = append(remoteAdded, remoteBranchName)
		} else if isRemoteRemoved, beforeRemoteBranchName, beforeRemoteBranchSHA := branchSpan.RemoteRemoved(); isRemoteRemoved {
			remoteRemoved[beforeRemoteBranchName] = beforeRemoteBranchSHA
		} else if isRemoteChanged, remoteBranchName, beforeSHA, afterSHA := branchSpan.RemoteChanged(); isRemoteChanged {
			remoteChanged[remoteBranchName] = undodomain.Change[gitdomain.SHA]{
				Before: beforeSHA,
				After:  afterSHA,
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
