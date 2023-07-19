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
		parent := l[current]
		if parent == "" {
			return result
		}
		result = append(Lineage{*parent}, result...)
		current = parent
	}
}

func (l Lineage) BranchNames() []string {
	result := make([]string, len(l))
	for b, branchInfo := range l {
		result[b] = branchInfo.Name
	}
	return result
}

// Children provides the names of all direct children of the given branch.
func (l Lineage) Children(branch string) []string {
	result := []string{}
	for child, parent := range l {
		if parent == branch {
			result = append(result, child)
		}
	}
	return result
}

// Contains indicates whether this Lineage contains a branch with the given name
func (l Lineage) Contains(branchName string) bool {
	_, ok := l[branchName]
	return ok
}

// HasParent returns whether or not the given branch has at least one parent.
func (l Lineage) HasParent(branch string) bool {
	parent, ok := l[branch]
	return parent != ""
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (l Lineage) IsAncestor(ancestor, other string) bool {
	current := other
	for {
		parent := l[current]
		if parent == "" {
			return false
		}
		if parent.Name == ancestor {
			return true
		}
		current = parent
	}
}

// OrderedHierarchically sorts the given BranchInfos so that ancestor branches come before their descendants
// and everything is sorted alphabetically.
func (l Lineage) OrderedHierarchically() Lineage {
	result := make(Lineage, len(l))
	copy(result, l)
	sort.Slice(result, func(x, y int) bool {
		return l.IsAncestor(result[x].Parent, result[y].Parent)
	})
	return result
}

// OrderedHierarchically provides the branches in this Lineage ordered so that ancestor branches come before their descendants
// and everything is sorted alphabetically.
func (l Lineage) OrderedHierarchically() []string {
	result := maps.Keys(l)
	sort.Slice(result, func(x, y int) bool {
		first := result[x]
		second := result[y]
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
	return result
}

// Roots provides the branches with children and no parents.
func (l Lineage) Roots() []string {
	roots := []string{}
	for _, parent := range l {
		_, found := l[parent]
		if !found && !stringslice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}
