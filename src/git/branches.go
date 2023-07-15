package git

import (
	"github.com/git-town/git-town/v9/src/config"
)

type BranchWithSyncStatus struct {
	Name       string
	SyncStatus SyncStatus
}

type BranchesWithSyncStatus []BranchWithSyncStatus

type BranchWithParentAndSyncStatus struct {
	Name       string
	Parent     string
	SyncStatus SyncStatus
}

// IsLocalBranch indicates whether this branch exists in the local repo that Git Town is running in.
func (bi BranchWithParentAndSyncStatus) IsLocal() bool {
	return bi.SyncStatus.IsLocal()
}

// Branches provides functionality for working with Git branches.
type BranchesWithParentAndSyncStatus []BranchWithParentAndSyncStatus

func NewBranchesWithParentAndSyncStatus(branchesWithParent []config.BranchWithParent, branchesWithSyncStatus BranchesWithSyncStatus) (BranchesWithParentAndSyncStatus, []string) {
	branchesWithParentAndSyncStatus := make(BranchesWithParentAndSyncStatus, len(branchesWithSyncStatus))
	for b, branchWithSyncStatus := range branchesWithSyncStatus {
		branchesWithParentAndSyncStatus[b] = BranchWithParentAndSyncStatus{
			Name:       branchWithSyncStatus.Name,
			SyncStatus: branchWithSyncStatus.SyncStatus,
		}
	}
	unused := []string{}
	for _, branchWithParent := range branchesWithParent {
		resultBranch := branchesWithParentAndSyncStatus.Lookup(branchWithParent.Name)
		if resultBranch == nil {
			unused = append(unused, branchWithParent.Name)
		} else {
			resultBranch.Parent = branchWithParent.Parent
		}
	}
	return branchesWithParentAndSyncStatus, unused
}

func (bs BranchesWithParentAndSyncStatus) BranchNames() []string {
	result := make([]string, len(bs))
	for b, branchInfo := range bs {
		result[b] = branchInfo.Name
	}
	return result
}

// LocalBranches provides only the branches that exist on the local machine.
func (bs BranchesWithParentAndSyncStatus) LocalBranches() BranchesWithParentAndSyncStatus {
	result := BranchesWithParentAndSyncStatus{}
	for _, branchInfo := range bs {
		if branchInfo.IsLocal() {
			result = append(result, branchInfo)
		}
	}
	return result
}

func (bs BranchesWithParentAndSyncStatus) Lookup(branch string) *BranchWithParentAndSyncStatus {
	for bi, branchInfo := range bs {
		if branchInfo.Name == branch {
			return &bs[bi]
		}
	}
	return nil
}
