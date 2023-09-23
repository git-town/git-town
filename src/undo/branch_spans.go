package undo

import "github.com/git-town/git-town/v9/src/domain"

type BranchSpans []BranchSpan

// Changes describes the changes made in this BranchesBeforeAfter structure.
func (bss BranchSpans) Changes() BranchChanges {
	result := EmptyBranchChanges()
	for _, branchSpan := range bss {
		if branchSpan.NoChanges() {
			continue
		}
		if branchSpan.IsOmniRemove() {
			result.OmniRemoved[branchSpan.Before.LocalName] = branchSpan.Before.LocalSHA
			continue
		}
		if branchSpan.IsOmniChange() {
			result.OmniChanged[branchSpan.Before.LocalName] = domain.Change[domain.SHA]{
				Before: branchSpan.Before.LocalSHA,
				After:  branchSpan.After.LocalSHA,
			}
			continue
		}
		if branchSpan.IsInconsistentChange() {
			result.InconsistentlyChanged = append(result.InconsistentlyChanged, domain.InconsistentChange{
				Before: branchSpan.Before,
				After:  branchSpan.After,
			})
			continue
		}
		switch {
		case branchSpan.LocalAdded():
			result.LocalAdded = append(result.LocalAdded, branchSpan.After.LocalName)
		case branchSpan.LocalRemoved():
			result.LocalRemoved[branchSpan.Before.LocalName] = branchSpan.Before.LocalSHA
		case branchSpan.LocalChanged():
			result.LocalChanged[branchSpan.Before.LocalName] = domain.Change[domain.SHA]{
				Before: branchSpan.Before.LocalSHA,
				After:  branchSpan.After.LocalSHA,
			}
		}
		switch {
		case branchSpan.RemoteAdded():
			result.RemoteAdded = append(result.RemoteAdded, branchSpan.After.RemoteName)
		case branchSpan.RemoteRemoved():
			result.RemoteRemoved[branchSpan.Before.RemoteName] = branchSpan.Before.RemoteSHA
		case branchSpan.RemoteChanged():
			result.RemoteChanged[branchSpan.Before.RemoteName] = domain.Change[domain.SHA]{
				Before: branchSpan.Before.RemoteSHA,
				After:  branchSpan.After.RemoteSHA,
			}
		}
	}
	return result
}
