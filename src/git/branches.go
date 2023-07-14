package git

import (
	"fmt"
	"sort"
)

// Branch contains the information Git Town needs to know about a Git branch.
type Branch struct {
	Name       string
	Parent     string
	SyncStatus SyncStatus
}

// IsLocalBranch indicates whether this branch exists in the local repo that Git Town is running in.
func (bi Branch) IsLocalBranch() bool {
	switch bi.SyncStatus {
	case SyncStatusLocalOnly, SyncStatusUpToDate, SyncStatusAhead, SyncStatusBehind:
		return true
	case SyncStatusRemoteOnly, SyncStatusDeletedAtRemote:
		return false
	}
	panic(fmt.Sprintf("uncaptured sync status: %v", bi.SyncStatus))
}

// Branches provides functionality for working with Git branches.
type Branches []Branch

// Ancestors provides all branches that are (great)(grand)parents of the branch with the given name.
func (bs Branches) Ancestors(branch string) Branches {
	result := Branches{}
	current := bs.Lookup(branch)
	if current == nil {
		return result
	}
	for {
		parent := bs.Lookup(current.Parent)
		if parent == nil {
			return result
		}
		result = append(Branches{*parent}, result...)
		current = parent
	}
}

func (bs Branches) BranchNames() []string {
	result := make([]string, len(bs))
	for b, branchInfo := range bs {
		result[b] = branchInfo.Name
	}
	return result
}

// Children provides all branches that have the branch with the given name as their parent.
func (bs Branches) Children(branchName string) Branches {
	result := Branches{}
	for _, b := range bs {
		if b.Parent == branchName {
			result = append(result, b)
		}
	}
	sort.Slice(result, func(a, b int) bool { return result[a].Name < result[b].Name })
	return result
}

// Contains indicates whether this Lineage contains a branch with the given name
func (bs Branches) Contains(branchName string) bool {
	for _, branch := range bs {
		if branch.Name == branchName {
			return true
		}
	}
	return false
}

// HasParents returns whether or not the given branch has at least one parent.
func (bs Branches) HasParent(branch string) bool {
	for _, branchInfo := range bs {
		if branchInfo.Parent == branch {
			return true
		}
	}
	return false
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (bs Branches) IsAncestor(ancestor, other string) bool {
	current := bs.Lookup(other)
	if current == nil {
		return false
	}
	for {
		parent := bs.Lookup(current.Parent)
		if parent == nil {
			return false
		}
		if parent.Name == ancestor {
			return true
		}
		current = parent
	}
}

// LocalBranches provides only the branches that exist on the local machine.
func (bs Branches) LocalBranches() Branches {
	result := Branches{}
	for _, branchInfo := range bs {
		if branchInfo.IsLocalBranch() {
			result = append(result, branchInfo)
		}
	}
	return result
}

func (bs Branches) Lookup(branch string) *Branch {
	for bi, branchInfo := range bs {
		if branchInfo.Name == branch {
			return &bs[bi]
		}
	}
	return nil
}

// OrderedHierarchically sorts the given BranchInfos so that ancestor branches come before their descendants
// and everything is sorted alphabetically.
func (bs Branches) OrderedHierarchically() Branches {
	result := make(Branches, len(bs))
	copy(result, bs)
	sort.Slice(result, func(x, y int) bool {
		return bs.IsAncestor(result[x].Parent, result[y].Parent)
	})
	return result
}

// Roots provides the branches without parents, i.e. branches that start a branch lineage.
func (bs Branches) Roots() Branches {
	roots := Branches{}
	for _, branch := range bs {
		if branch.Parent == "" && !roots.Contains(branch.Name) {
			roots = append(roots, branch)
		}
	}
	sort.Slice(roots, func(a, b int) bool {
		return roots[a].Name < roots[b].Name
	})
	return roots
}
