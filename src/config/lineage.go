package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/git-town/git-town/v9/src/stringslice"
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
		parent := l.Lookup(current.Parent)
		fmt.Printf("LOOKING UP PARENT BRANCH %q\n", current.Parent)
		fmt.Printf("LINEAGE IS %q\n", strings.Join(l.BranchNames(), ", "))
		entries = append([]BranchWithParent{*parent}, entries...)
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
		parent := l.Lookup(current.Parent)
		if parent == nil {
			return false
		}
		if parent.Name == ancestorName {
			return true
		}
		current = parent
	}
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
	sort.Slice(roots, func(a, b int) bool {
		return roots[a].Name < roots[b].Name
	})
	return roots
}

// OrderedHierarchically sorts the given BranchInfos so that ancestor branches come before their descendants
// and everything is sorted alphabetically.
func (l Lineage) OrderedHierarchically() Lineage {
	result := make(Lineage, len(l))
	copy(result, l)
	sort.Slice(result, func(x, y int) bool {
		return l.IsAncestor(result[x].Name, result[y].Name)
		// fmt.Printf("%s %s %v", result[x].Name)
	})
	return result
}
