package undobranches

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/undo/undodomain"
)

// BranchSpans describes how a Git Town command has modified the branches in a Git repository.
type BranchSpans []BranchSpan

func NewBranchSpans(beforeSnapshot, afterSnapshot gitdomain.BranchesSnapshot) BranchSpans {
	result := BranchSpans{}
	for _, before := range beforeSnapshot.Branches {
		after := afterSnapshot.Branches.FindMatchingRecord(before)
		result = append(result, BranchSpan{
			Before: before,
			After:  after,
		})
	}
	for _, after := range afterSnapshot.Branches {
		if beforeSnapshot.Branches.FindMatchingRecord(after).IsEmpty() {
			result = append(result, BranchSpan{
				Before: gitdomain.EmptyBranchInfo(),
				After:  after,
			})
		}
	}
	return result
}

// Changes describes the specific changes made in this BranchSpans.
func (self BranchSpans) Changes() BranchChanges {
	result := EmptyBranchChanges()
	for _, branchSpan := range self {
		if branchSpan.NoChanges() {
			continue
		}
		if isOmniRemove, beforeLocalBranch, beforeLocalSHA := branchSpan.IsOmniRemove2(); isOmniRemove {
			result.OmniRemoved[beforeLocalBranch] = beforeLocalSHA
			continue
		}
		if isOmniChange, branchName, beforeSHA, afterSHA := branchSpan.IsOmniChange2(); isOmniChange {
			result.OmniChanged[branchName] = undodomain.Change[gitdomain.SHA]{
				Before: beforeSHA,
				After:  afterSHA,
			}
			continue
		}
		if branchSpan.IsInconsistentChange() {
			result.InconsistentlyChanged = append(result.InconsistentlyChanged, undodomain.InconsistentChange{
				Before: branchSpan.Before,
				After:  branchSpan.After,
			})
			continue
		}
		if localAdded, afterBranch, _ := branchSpan.LocalAdded2(); localAdded {
			result.LocalAdded = append(result.LocalAdded, afterBranch)
		} else if localRemoved, beforeBranch, beforeSHA := branchSpan.LocalRemoved2(); localRemoved {
			result.LocalRemoved[beforeBranch] = beforeSHA
		} else if localChanged, branch, beforeSHA, afterSHA := branchSpan.LocalChanged2(); localChanged {
			result.LocalChanged[branch] = undodomain.Change[gitdomain.SHA]{
				Before: beforeSHA,
				After:  afterSHA,
			}
		}
		switch {
		case branchSpan.RemoteAdded():
			result.RemoteAdded = append(result.RemoteAdded, branchSpan.After.RemoteName)
		case branchSpan.RemoteRemoved():
			result.RemoteRemoved[branchSpan.Before.RemoteName] = branchSpan.Before.RemoteSHA
		case branchSpan.RemoteChanged():
			result.RemoteChanged[branchSpan.Before.RemoteName] = undodomain.Change[gitdomain.SHA]{
				Before: branchSpan.Before.RemoteSHA,
				After:  branchSpan.After.RemoteSHA,
			}
		}
	}
	return result
}
