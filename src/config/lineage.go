package config

import (
	"sort"

	"github.com/git-town/git-town/v9/src/stringslice"
)

// Lineage encapsulates all data and functionality around parent branches.
type Lineage struct {
	// branch --> its parent
	Entries map[string]string
}

// Ancestors provides the names of all parent branches of the branch with the given name.
func (l *Lineage) Ancestors(branch string) []string {
	current := branch
	result := []string{}
	for {
		parent, found := l.Entries[current]
		if !found {
			return result
		}
		result = append([]string{parent}, result...)
		current = parent
	}
}

// Children provides the names of all branches that have the given branch as their parent.
func (l *Lineage) Children(branch string) []string {
	result := []string{}
	for child, parent := range l.Entries {
		if parent == branch {
			result = append(result, child)
		}
	}
	sort.Strings(result)
	return result
}

// HasParents returns whether or not the given branch has at least one parent.
func (l *Lineage) HasParents(branch string) bool {
	for child := range l.Entries {
		if child == branch {
			return true
		}
	}
	return false
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (l *Lineage) IsAncestor(ancestor, other string) bool {
	current := other
	for {
		parent, found := l.Entries[current]
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
func (l *Lineage) Parent(branch string) string {
	for child, parent := range l.Entries {
		if child == branch {
			return parent
		}
	}
	return ""
}

// Roots provides the branches with children and no parents.
func (l *Lineage) Roots() []string {
	roots := []string{}
	for _, parent := range l.Entries {
		_, ok := l.Entries[parent]
		if !ok && !stringslice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}

func (l *Lineage) SetParent(branch, parent string) {
	l.Entries[branch] = parent
}
