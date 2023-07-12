package git

import (
	"fmt"
	"sort"

	"github.com/git-town/git-town/v9/src/config"
)

// BranchInfo contains information about the sync status of a Git branch.
type BranchInfo struct {
	Name       string
	Parent     string
	SyncStatus SyncStatus
}

func (bi BranchInfo) IsLocalBranch() bool {
	switch bi.SyncStatus {
	case SyncStatusLocalOnly, SyncStatusUpToDate, SyncStatusAhead, SyncStatusBehind:
		return true
	case SyncStatusRemoteOnly, SyncStatusDeletedAtRemote:
		return false
	}
	panic(fmt.Sprintf("uncaptured sync status: %v", bi.SyncStatus))
}

type BranchInfos []BranchInfo

// LocalBranches provides only the branches that exist on the local machine.
func (bi BranchInfos) LocalBranches() BranchInfos {
	result := BranchInfos{}
	for _, branchInfo := range bi {
		if branchInfo.IsLocalBranch() {
			result = append(result, branchInfo)
		}
	}
	return result
}

func (bi BranchInfos) Lookup(branch string) *BranchInfo {
	for b, branchInfo := range bi {
		if branchInfo.Name == branch {
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

// OrderedHierarchically sorts the given BranchInfos so that ancestor branches come before their descendants
// and everything is sorted alphabetically.
func (bi BranchInfos) OrderedHierarchically() BranchInfos {
	result := make(BranchInfos, len(bi))
	copy(result, bi)
	lineage := bi.lineage()
	sort.Slice(result, func(a, b int) bool {
		return lineage.IsAncestor(result[a].Parent, result[b].Parent)
	})
	return result
}

// provides a Lineage instance for these branches
func (bi BranchInfos) lineage() config.Lineage {
	parents := map[string]string{}
	for _, branchInfo := range bi {
		parents[branchInfo.Name] = branchInfo.Parent
	}
	return config.Lineage{Entries: parents}
}
