package undobranches

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/undo/undodomain"
)

// BranchSpans describes how a Git Town command has modified the branches in a Git repository.
type BranchSpans []BranchSpan

func NewBranchSpans(beforeSnapshot, afterSnapshot gitdomain.BranchesSnapshot) BranchSpans {
	result := BranchSpans{}
	for _, before := range beforeSnapshot.Branches {
		after := afterSnapshot.Branches.FindMatchingRecord(before)
		result = append(result, BranchSpan{
			Before: Some(before),
			After:  after,
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
	omniChanged := LocalBranchChange{}
	omniRemoved := LocalBranchesSHAs{}
	remoteAdded := gitdomain.RemoteBranchNames{}
	remoteChanged := map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{}
	remoteRemoved := map[gitdomain.RemoteBranchName]gitdomain.SHA{}
	for _, branchSpan := range self {
		if branchSpan.NoChanges() {
			continue
		}
		if isOmniRemove, beforeLocalBranch, beforeLocalSHA := branchSpan.IsOmniRemove(); isOmniRemove {
			omniRemoved[beforeLocalBranch] = beforeLocalSHA
			continue
		}
		if isOmniChange, branchName, beforeSHA, afterSHA := branchSpan.IsOmniChange(); isOmniChange {
			omniChanged[branchName] = undodomain.Change[gitdomain.SHA]{
				Before: beforeSHA,
				After:  afterSHA,
			}
			continue
		}
		isInconsistentChange, before, after := branchSpan.IsInconsistentChange()
		if isInconsistentChange {
			inconsistentlyChanged = append(inconsistentlyChanged, undodomain.InconsistentChange{
				Before: before,
				After:  after,
			})
			continue
		}
		if isLocalAdded, afterBranch, _ := branchSpan.LocalAdded(); isLocalAdded {
			localAdded = append(localAdded, afterBranch)
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
		OmniChanged:           omniChanged,
		OmniRemoved:           omniRemoved,
		RemoteAdded:           remoteAdded,
		RemoteChanged:         remoteChanged,
		RemoteRemoved:         remoteRemoved,
	}
}
