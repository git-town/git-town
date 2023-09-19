package undo

import "github.com/git-town/git-town/v9/src/domain"

type BranchSpans []BranchSpan

// Changes describes the changes made in this BranchesBeforeAfter structure.
func (bs BranchSpans) Changes() Changes {
	result := EmptyChanges()
	for _, ba := range bs {
		if ba.NoChanges() {
			continue
		}
		if ba.IsOmniChange() {
			result.OmniChanged[ba.Before.LocalName] = domain.Change[domain.SHA]{
				Before: ba.Before.LocalSHA,
				After:  ba.After.LocalSHA,
			}
			continue
		}
		if ba.IsInconsintentChange() {
			result.InconsistentlyChanged = append(result.InconsistentlyChanged, domain.InconsistentChange{
				Before: ba.Before,
				After:  ba.After,
			})
			continue
		}
		switch {
		case ba.LocalAdded():
			result.LocalAdded = append(result.LocalAdded, ba.After.LocalName)
		case ba.LocalRemoved():
			result.LocalRemoved[ba.Before.LocalName] = ba.Before.LocalSHA
		case ba.LocalChanged():
			result.LocalChanged[ba.Before.LocalName] = domain.Change[domain.SHA]{
				Before: ba.Before.LocalSHA,
				After:  ba.After.LocalSHA,
			}
		}
		switch {
		case ba.RemoteAdded():
			result.RemoteAdded = append(result.RemoteAdded, ba.After.RemoteName)
		case ba.RemoteRemoved():
			result.RemoteRemoved[ba.Before.RemoteName] = ba.Before.RemoteSHA
		case ba.RemoteChanged():
			result.RemoteChanged[ba.Before.RemoteName] = domain.Change[domain.SHA]{
				Before: ba.Before.RemoteSHA,
				After:  ba.After.RemoteSHA,
			}
		}
	}
	return result
}
