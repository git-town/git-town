package git

type BranchWithSyncStatus struct {
	Name       string
	SyncStatus SyncStatus
}

type BranchesWithSyncStatus []BranchWithSyncStatus

func (bs BranchesWithSyncStatus) BranchNames() []string {
	result := make([]string, len(bs))
	for b, branchInfo := range bs {
		result[b] = branchInfo.Name
	}
	return result
}

// IsLocalBranch indicates whether this branch exists in the local repo that Git Town is running in.
func (bi BranchWithSyncStatus) IsLocal() bool {
	return bi.SyncStatus.IsLocal()
}

// LocalBranches provides only the branches that exist on the local machine.
func (bs BranchesWithSyncStatus) LocalBranches() BranchesWithSyncStatus {
	result := BranchesWithSyncStatus{}
	for _, branchInfo := range bs {
		if branchInfo.IsLocal() {
			result = append(result, branchInfo)
		}
	}
	return result
}

func (bs BranchesWithSyncStatus) Lookup(branch string) *BranchWithSyncStatus {
	for bi, branchInfo := range bs {
		if branchInfo.Name == branch {
			return &bs[bi]
		}
	}
	return nil
}
