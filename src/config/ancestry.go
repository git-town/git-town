package config

import (
	"sort"
	"strings"

	"github.com/git-town/git-town/v9/src/stringslice"
)

// Ancestry encapsulates all data and functionality around parent branches.
type Ancestry struct {
	// branch --> its parent
	entries    map[string]string
	mainBranch string
}

func NewAncestry(mainBranch string) Ancestry {
	return Ancestry{map[string]string{}, mainBranch}
}

// ParentBranchMap returns a map from branch name to its parent branch.
func LoadAncestry(git Git, mainBranch string) Ancestry {
	parents := map[string]string{}
	for _, key := range git.LocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := git.LocalConfigValue(key)
		parents[child] = parent
	}
	return Ancestry{parents, mainBranch}
}

func (a *Ancestry) AddParent(branch, parent string) {
	a.entries[branch] = parent
}

// Ancestors provides the names of all parent branches for the given branch,
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
// TODO: rename to Ancestors
func (a *Ancestry) Ancestors(branch string) []string {
	current := branch
	result := []string{}
	for {
		parent, found := a.entries[current]
		if !found {
			return result
		}
		result = append([]string{parent}, result...)
		current = parent
	}
}

// Children provides the names of all branches that have the given branch as their parent.
func (a *Ancestry) Children(branch string) []string {
	result := []string{}
	for child, parent := range a.entries {
		if parent == branch {
			result = append(result, child)
		}
	}
	sort.Strings(result)
	return result
}

// branch --> its parent
func (a *Ancestry) Entries() map[string]string {
	return a.entries
}

// HasParent returns whether or not the given branch has a parent.
func (a *Ancestry) HasParent(branch string) bool {
	for child := range a.entries {
		if child == branch {
			return true
		}
	}
	return false
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (a *Ancestry) IsAncestor(branch, ancestor string) bool {
	ancestors := a.Ancestors(branch)
	return stringslice.Contains(ancestors, ancestor)
}

// Parent provides the name of the parent branch for the given branch or nil if the branch has no parent.
func (a *Ancestry) Parent(branch string) string {
	for child, parent := range a.entries {
		if child == branch {
			return parent
		}
	}
	return ""
}

// Roots provides the branches with children and no parents.
func (a *Ancestry) Roots() []string {
	roots := []string{}
	for _, parent := range a.entries {
		_, ok := a.entries[parent]
		if !ok && !stringslice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}
