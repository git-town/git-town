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
		if remoteAddedResult := branchSpan.RemoteAdded(); remoteAddedResult.IsAdded {
			remoteAdded = append(remoteAdded, remoteAddedResult.AddedRemoteBranchName)
		} else if remoteRemovedResult := branchSpan.RemoteRemoved(); remoteRemovedResult.IsRemoved {
			remoteRemoved[remoteRemovedResult.RemoteBranchName] = remoteRemovedResult.BeforeRemoteSHA
		} else if remoteChangedResult := branchSpan.RemoteChanged(); remoteChangedResult.IsChanged {
			remoteChanged[remoteAddedResult.AddedRemoteBranchName] = undodomain.Change[gitdomain.SHA]{
				Before: remoteChangedResult.BeforeSHA,
				After:  remoteChangedResult.AfterSHA,
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
