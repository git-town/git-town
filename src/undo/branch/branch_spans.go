package branch

import "github.com/git-town/git-town/v9/src/domain"

// Spans describes how a Git Town command has modified the branches in a Git repository.
type Spans []Span

func NewSpans(beforeSnapshot, afterSnapshot domain.BranchesSnapshot) Spans {
	result := Spans{}
	for _, before := range beforeSnapshot.Branches {
		after := afterSnapshot.Branches.FindMatchingRecord(before)
		result = append(result, Span{
			Before: before,
			After:  after,
		})
	}
	for _, after := range afterSnapshot.Branches {
		if beforeSnapshot.Branches.FindMatchingRecord(after).IsEmpty() {
			result = append(result, Span{
				Before: domain.EmptyBranchInfo(),
				After:  after,
			})
		}
	}
	return result
}

// Changes describes the specific changes made in this BranchSpans.
func (self Spans) Changes() Changes {
	result := EmptyChanges()
	for _, branchSpan := range self {
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
