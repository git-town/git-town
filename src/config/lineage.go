package config

import (
	"sort"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/slice"
	"golang.org/x/exp/maps"
)

// Lineage encapsulates all data and functionality around parent branches.
// branch --> its parent.
type Lineage map[domain.LocalBranchName]domain.LocalBranchName

// BranchAndAncestors provides the full lineage for the branch with the given name,
// including the branch.
func (l Lineage) BranchAndAncestors(branchName domain.LocalBranchName) domain.LocalBranchNames {
	return append(l.Ancestors(branchName), branchName)
}

// BranchesAndAncestors provides the full lineage for the branches with the given names,
// including the branches themselves.
func (l Lineage) BranchesAndAncestors(branchNames domain.LocalBranchNames) domain.LocalBranchNames {
	result := branchNames
	for _, branchName := range branchNames {
		ancestors := l.Ancestors(branchName)
		result = slice.AppendAllMissing(result, ancestors)
	}
	l.OrderHierarchically(result)
	return result
}

// Ancestors provides the names of all parent branches of the branch with the given name.
func (l Lineage) Ancestors(branch domain.LocalBranchName) domain.LocalBranchNames {
	current := branch
	result := domain.LocalBranchNames{}
	for {
		parent, found := l[current]
		if !found {
			return result
		}
		result = append(domain.LocalBranchNames{parent}, result...)
		current = parent
	}
}

// BranchNames provides the names of all branches in this Lineage, sorted alphabetically.
func (l Lineage) BranchNames() domain.LocalBranchNames {
	result := domain.LocalBranchNames(maps.Keys(l))
	result.Sort()
	return result
}

// Children provides the names of all branches that have the given branch as their parent.
func (l Lineage) Children(branch domain.LocalBranchName) domain.LocalBranchNames {
	result := domain.LocalBranchNames{}
	for child, parent := range l {
		if parent == branch {
			result = append(result, child)
		}
	}
	result.Sort()
	return result
}

// HasParents returns whether or not the given branch has at least one parent.
func (l Lineage) HasParents(branch domain.LocalBranchName) bool {
	for child := range l {
		if child == branch {
			return true
		}
	}
	return false
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (l Lineage) IsAncestor(ancestor, other domain.LocalBranchName) bool {
	current := other
	for {
		parent, found := l[current]
		if !found {
			return false
		}
		if parent == ancestor {
			return true
		}
		current = parent
	}
}

// Parent provides the name of the parent branch for the given branch or nil if the branch has no parent.
func (l Lineage) Parent(branch domain.LocalBranchName) domain.LocalBranchName {
	for child, parent := range l {
		if child == branch {
			return parent
		}
	}
	return domain.LocalBranchName{}
}

// OrderHierarchically sorts the given branches so that ancestor branches come before their descendants
// and everything is sorted alphabetically.
func (l Lineage) OrderHierarchically(branches domain.LocalBranchNames) {
	sort.Slice(branches, func(a, b int) bool {
		first := branches[a]
		second := branches[b]
		if first.IsEmpty() {
			return true
		}
		if second.IsEmpty() {
			return false
		}
		isAncestor := l.IsAncestor(first, second)
		if isAncestor {
			return true
		}
		return first.String() < second.String()
	})
}

// Roots provides the branches with children and no parents.
func (l Lineage) Roots() domain.LocalBranchNames {
	roots := domain.LocalBranchNames{}
	for _, parent := range l {
		_, found := l[parent]
		if !found && !slice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	roots.Sort()
	return roots
}
