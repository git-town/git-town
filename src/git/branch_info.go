package git

import (
	"fmt"
	"sort"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
)

// BranchInfo contains information about the sync status of a Git branch.
type BranchInfo struct {
	Name       string
	Parent     string
	SyncStatus SyncStatus
}

type BranchInfos []BranchInfo

// LocalBranches provides only the branches that exist on the local machine.
func (bi BranchInfos) LocalBranches() BranchInfos {
	result := BranchInfos{}
	for _, branchInfo := range bi {
		var isLocalBranch bool
		switch branchInfo.SyncStatus {
		case SyncStatusLocalOnly, SyncStatusUpToDate, SyncStatusAhead, SyncStatusBehind:
			isLocalBranch = true
		case SyncStatusRemoteOnly, SyncStatusDeletedAtRemote:
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

// OrderedHierarchically sorts the given BranchInfos so that ancestor branches come before their descendants
// and everything is sorted alphabetically.
func (bi BranchInfos) OrderedHierarchically() []BranchInfo {
	// for now we just put the main branch first
	// TODO: sort this better by putting parent branches before child branches
	result := make([]BranchInfo, len(bi))
	copy(result, bi)
	sort.Slice(result, func(a, b int) bool {
		ap := result[a].Parent
		bp := result[b].Parent
		fmt.Printf("COMPARING %s with %s ...", ap, bp)
		if ap == "" {
			fmt.Println("true")
			return true
		}
		if bp == "" {
			fmt.Println("false")
			return false
		}
		if ap == mainBranch {
			fmt.Println("true")
			return true
		}
		if bp == mainBranch {
			fmt.Println("false")
			return false
		}
		result := ap < bp
		fmt.Printf("%v\n", result)
		return result
	})
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

// provides a Lineage instance for these branches
func (bi BranchInfos) lineage() config.Lineage {
	parents := map[string]string{}
	for _, key := range gt.LocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := gt.LocalConfigValue(key)
		parents[child] = parent
	}
	return config.Lineage{parents, gt.MainBranch()}
}
