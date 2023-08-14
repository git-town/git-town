package config

import (
	"sort"

	"github.com/git-town/git-town/v9/src/slice"
	"golang.org/x/exp/maps"
)

// Lineage encapsulates all data and functionality around parent branches.
// branch --> its parent.
type Lineage map[string]string

// BranchAndAncestors provides the full lineage for the branch with the given name,
// including the branch.
func (l Lineage) BranchAndAncestors(branchName string) []string {
	return append(l.Ancestors(branchName), branchName)
}

// BranchesAndAncestors provides the full lineage for the branches with the given names,
// including the branches themselves.
func (l Lineage) BranchesAndAncestors(branchNames []string) []string {
	result := branchNames
	for _, branchName := range branchNames {
		ancestors := l.Ancestors(branchName)
		result = slice.AppendAllMissing(result, ancestors)
	}
	l.OrderHierarchically(result)
	return result
}

// Ancestors provides the names of all parent branches of the branch with the given name.
func (l Lineage) Ancestors(branch string) []string {
	current := branch
	result := []string{}
	for {
		parent := l[current]
		if parent == "" {
			return result
		}
		result = append([]string{parent}, result...)
		current = parent
	}
}

// BranchNames provides the names of all branches in this Lineage, sorted alphabetically.
func (l Lineage) BranchNames() []string {
	result := maps.Keys(l)
	sort.Strings(result)
	return result
}

// Children provides the names of all branches that have the given branch as their parent.
func (l Lineage) Children(branch string) []string {
	result := []string{}
	for child, parent := range l {
		if parent == branch {
			result = append(result, child)
		}
	}
	sort.Strings(result)
	return result
}

// HasParents returns whether or not the given branch has at least one parent.
func (l Lineage) HasParents(branch string) bool {
	for child := range l {
		if child == branch {
			return true
		}
	}
	return false
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (l Lineage) IsAncestor(ancestor, other string) bool {
	current := other
	for {
		parent := l[current]
		if parent == "" {
			return false
		}
		if parent == ancestor {
			return true
		}
		current = parent
	}
}

// Parent provides the name of the parent branch for the given branch or nil if the branch has no parent.
func (l Lineage) Parent(branch string) string {
	for child, parent := range l {
		if child == branch {
			return parent
		}
	}
	return ""
}

// OrderHierarchically sorts the given branches so that ancestor branches come before their descendants
// and everything is sorted alphabetically.
func (l Lineage) OrderHierarchically(branches []string) {
	sort.Slice(branches, func(x, y int) bool {
		first := branches[x]
		second := branches[y]
		if first == "" {
			return true
		}
		if second == "" {
			return false
		}
		isAncestor := l.IsAncestor(first, second)
		if isAncestor {
			return true
		}
		return first < second
	})
}

// Roots provides the branches with children and no parents.
func (l Lineage) Roots() []string {
	roots := []string{}
	for _, parent := range l {
		_, found := l[parent]
		if !found && !slice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}
