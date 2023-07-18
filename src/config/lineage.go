package config

import (
	"sort"

	"github.com/git-town/git-town/v9/src/stringslice"
	"golang.org/x/exp/maps"
)

// Lineage encapsulates all data and functionality around parent branches.
// branch --> its parent.
type Lineage map[string]string

// Ancestors provides the names of all parent branches of the branch with the given name.
func (l Lineage) Ancestors(branch string) []string {
	current := branch
	result := []string{}
	for {
		parent, found := l[current]
		if !found {
			return result
		}
		result = append([]string{parent}, result...)
		current = parent
	}
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
func (l Lineage) Parent(branch string) string {
	for child, parent := range l {
		if child == branch {
			return parent
		}
	}
	return ""
}

// Roots provides the branches with children and no parents.
func (l Lineage) Roots() []string {
	roots := []string{}
	for _, parent := range l {
		_, ok := l[parent]
		if !ok && !stringslice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}

// OrderedHierarchically provides the branches in this Lineage ordered so that ancestor branches come before their descendants
// and everything is sorted alphabetically.
func (l Lineage) OrderedHierarchically() []string {
	result := maps.Keys(l)
	sort.Slice(result, func(x, y int) bool {
		return l.IsAncestor(result[x], result[y])
	})
	return result
}
