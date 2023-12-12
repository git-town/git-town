package domain

import (
	"sort"
	"strings"
)

type LocalBranchNames []LocalBranchName

func NewLocalBranchNames(names ...string) LocalBranchNames {
	result := make(LocalBranchNames, len(names))
	for n, name := range names {
		result[n] = NewLocalBranchName(name)
	}
	return result
}

func (self LocalBranchNames) Categorize(branchTypes BranchTypes) (perennials, features LocalBranchNames) {
	for _, branch := range self {
		if branchTypes.IsFeatureBranch(branch) {
			features = append(features, branch)
		} else {
			perennials = append(perennials, branch)
		}
	}
	return
}

// Hoist moves the given needle to the front of the list.
func (self LocalBranchNames) Hoist(needle LocalBranchName) LocalBranchNames {
	result := make(LocalBranchNames, 0, len(self))
	foundNeedle := false
	for _, branch := range self {
		if branch == needle {
			foundNeedle = true
		} else {
			result = append(result, branch)
		}
	}
	if foundNeedle {
		result = append(LocalBranchNames{needle}, result...)
	}
	return result
}

// Join provides the names of all branches in this collection connected by the given separator.
func (self LocalBranchNames) Join(sep string) string {
	return strings.Join(self.Strings(), sep)
}

// Remove removes the given branch name from this collection.
func (self LocalBranchNames) Remove(toRemove LocalBranchName) LocalBranchNames {
	result := make(LocalBranchNames, 0, len(self)-1)
	for _, branch := range self {
		if branch != toRemove {
			result = append(result, branch)
		}
	}
	return result
}

// RemoveWorkspaceMarkers removes the workspace markers from the branch names in this list.
func (self LocalBranchNames) RemoveWorkspaceMarkers() LocalBranchNames {
	result := make(LocalBranchNames, len(self))
	for b, branch := range self {
		if strings.HasPrefix(branch.String(), "+ ") {
			result[b] = branch[2:]
		} else {
			result[b] = branch
		}
	}
	return result
}

// Sort orders the branches in this collection alphabetically.
func (self LocalBranchNames) Sort() {
	sort.Slice(self, func(a, b int) bool {
		return self[a] < self[b]
	})
}

// Strings provides the names of all branches in this collection as strings.
func (self LocalBranchNames) Strings() []string {
	result := make([]string, len(self))
	for b, branch := range self {
		result[b] = branch.String()
	}
	return result
}
