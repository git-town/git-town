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
	for b, branch := range bs {
		result[b] = branch.Name
	}
	return result
}

func (bs BranchesSyncStatus) Contains(branchName string) bool {
	for _, branch := range bs {
		if branch.Name == branchName {
			return true
		}
	}
	return false
}

// LocalBranches provides only the branches that exist on the local machine.
func (bs BranchesSyncStatus) LocalBranches() BranchesSyncStatus {
	result := BranchesSyncStatus{}
	for _, branch := range bs {
		if branch.IsLocal() {
			result = append(result, branch)
		}
	}
	return result
}

func (bs BranchesSyncStatus) Lookup(branchName string) *BranchSyncStatus {
	for bi, branch := range bs {
		if branch.Name == branchName {
			return &bs[bi]
		}
	}
	return nil
}
