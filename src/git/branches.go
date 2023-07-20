package git

type BranchSyncStatus struct {
	Name       string
	SyncStatus SyncStatus
}

// IsLocalBranch indicates whether this branch exists in the local repo that Git Town is running in.
func (bi BranchSyncStatus) IsLocal() bool {
	return bi.SyncStatus.IsLocal()
}

type BranchesSyncStatus []BranchSyncStatus

func (bs BranchesSyncStatus) BranchNames() []string {
	result := make([]string, len(bs))
	for b, branchInfo := range bs {
		result[b] = branchInfo.Name
	}
	return result
}

// LocalBranches provides only the branches that exist on the local machine.
func (bs BranchesSyncStatus) LocalBranches() BranchesSyncStatus {
	result := BranchesSyncStatus{}
	for _, branchInfo := range bs {
		if branchInfo.IsLocal() {
			result = append(result, branchInfo)
		}
	}
	return result
}

func (bs BranchesSyncStatus) Lookup(branchName string) *BranchSyncStatus {
	for bi, branchInfo := range bs {
		if branchInfo.Name == branchName {
			return &bs[bi]
		}
	}
	return nil
}
