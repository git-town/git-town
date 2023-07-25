package git

import "fmt"

type BranchSyncStatus struct {
	Name       string
	SyncStatus SyncStatus
}

func (bi BranchSyncStatus) HasTrackingBranch() bool {
	switch bi.SyncStatus {
	case SyncStatusAhead, SyncStatusBehind, SyncStatusUpToDate, SyncStatusRemoteOnly:
		return true
	case SyncStatusLocalOnly, SyncStatusDeletedAtRemote:
		return false
	}
	panic(fmt.Sprintf("unknown sync status: %v", bi.SyncStatus))
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

// LocalBranchesWithDeletedTrackingBranches provides only the branches that exist locally and have a deleted tracking branch.
func (bs BranchesSyncStatus) LocalBranchesWithDeletedTrackingBranches() BranchesSyncStatus {
	result := BranchesSyncStatus{}
	for _, branch := range bs {
		if branch.SyncStatus == SyncStatusDeletedAtRemote {
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

func (bs BranchesSyncStatus) Select(names []string) (BranchesSyncStatus, error) {
	result := make(BranchesSyncStatus, len(names))
	for n, name := range names {
		branch := bs.Lookup(name)
		if branch == nil {
			return result, fmt.Errorf("cannot find branch %q", name)
		}
		result[n] = *branch
	}
	return result, nil
}
