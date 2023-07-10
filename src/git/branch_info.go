package git

// BranchInfo contains information about the sync status of a Git branch.
type BranchInfo struct {
	Name       string
	Location   BranchLocation
	SyncStatus SyncStatus
}

type BranchInfos []BranchInfo

// LocalBranches provides only the branches that exist on the local machine.
func (bi BranchInfos) LocalBranches() BranchInfos {
	result := BranchInfos{}
	for _, branchInfo := range bi {
		var isLocalBranch bool
		switch branchInfo.Location {
		case BranchLocationLocalOnly, BranchLocationLocalAndRemote, BranchLocationDeletedAtRemote:
			isLocalBranch = true
		case BranchLocationRemoteOnly:
			isLocalBranch = false
		}
		if isLocalBranch {
			result = append(result, branchInfo)
		}
	}
	return result
}

func (bi BranchInfos) Lookup(branch string) *BranchInfo {
	for b := range bi {
		if bi[b].Name == branch {
			return &bi[b]
		}
	}
	return nil
}

func (bi BranchInfos) BranchNames() []string {
	result := make([]string, len(bi))
	for b, branchInfo := range bi {
		result[b] = branchInfo.Name
	}
	return result
}

// OrderedHierarchically sorts the given BranchInfo list so that ancestor branches come before their descendants and everything is sorted alphabetically.
func (bi BranchInfos) OrderedHierarchically() []BranchInfo {
	// for now we just put the main branch first
	// TODO: sort this better by putting parent branches before child branches
	result := make([]BranchInfo, len(bi))
	// for b, branchInfo := range bi {
	// 	// result[b] =
	// }
	return result
}

// IndexOfBranch returns the zero-based index of the branch with the given name.
func (bi BranchInfos) IndexOfBranch(branch string) (pos int, found bool) {
	for b, branchInfo := range bi {
		if branchInfo.Name == branch {
			return b, true
		}
	}
	return 0, false
}
