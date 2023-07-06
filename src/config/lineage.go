package config

import (
	"sort"
	"strings"

	"github.com/git-town/git-town/v9/src/stringslice"
)

// Lineage encapsulates all data and functionality around parent branches.
type Lineage struct {
	// branch --> its parent
	entries    map[string]string
	mainBranch string
}

func NewLineage(mainBranch string) Lineage {
	return Lineage{map[string]string{}, mainBranch}
}

// LoadLineage provides the Lineage for the given Git repo.
func LoadLineage(git Git, mainBranch string) Lineage {
	parents := map[string]string{}
	for _, key := range git.LocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := git.LocalConfigValue(key)
		parents[child] = parent
	}
	return Lineage{parents, mainBranch}
}

func (l *Lineage) AddParent(branch, parent string) {
	l.entries[branch] = parent
}

// Ancestors provides the names of all parent branches for the given branch,
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
// TODO: rename to Ancestors
func (l *Lineage) Ancestors(branch string) []string {
	current := branch
	result := []string{}
	for {
		parent, found := l.entries[current]
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
	for child, parent := range l.entries {
		if parent == branch {
			result = append(result, child)
		}
	}
	sort.Strings(result)
	return result
}

// branch --> its parent
func (l *Lineage) Entries() map[string]string {
	return l.entries
}

// HasParents returns whether or not the given branch has at least one parent.
func (a *Lineage) HasParents(branch string) bool {
	for child := range a.entries {
		if child == branch {
			return true
		}
	}
	return false
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (l *Lineage) IsAncestor(branch, ancestor string) bool {
	ancestors := l.Ancestors(branch)
	return stringslice.Contains(ancestors, ancestor)
}

// Parent provides the name of the parent branch for the given branch or nil if the branch has no parent.
func (l *Lineage) Parent(branch string) string {
	for child, parent := range l.entries {
		if child == branch {
			return parent
		}
	}
	return ""
}

// Roots provides the branches with children and no parents.
func (l *Lineage) Roots() []string {
	roots := []string{}
	for _, parent := range l.entries {
		_, ok := l.entries[parent]
		if !ok && !stringslice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}
