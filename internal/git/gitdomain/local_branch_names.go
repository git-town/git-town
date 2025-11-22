package gitdomain

import (
	"slices"
	"strings"

	"github.com/git-town/git-town/v22/internal/gohacks/slice"
)

type LocalBranchNames []LocalBranchName

func NewLocalBranchNames(names ...string) LocalBranchNames {
	result := make(LocalBranchNames, len(names))
	for n, name := range names {
		result[n] = NewLocalBranchName(name)
	}
	return result
}

// ParseLocalBranchNames constructs a LocalBranchNames instance
// containing the branches listed in the given space-separated string.
func ParseLocalBranchNames(names string) LocalBranchNames {
	parts := strings.Split(names, " ")
	result := make(LocalBranchNames, 0, len(parts))
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		result = append(result, NewLocalBranchName(part))
	}
	return result
}

// Intersect provides all branch names that occur in this list and the given list,
// ordered the same as in this list.
func (self LocalBranchNames) Intersect(branches LocalBranchNames) LocalBranchNames {
	result := make(LocalBranchNames, 0, len(self))
	for _, branch := range self {
		if branches.Contains(branch) {
			result = append(result, branch)

		}
	}
	return result
}

// AppendAllMissing provides a LocalBranchNames list consisting of the sum of this and elements of other list that aren't in this list.
func (self LocalBranchNames) AppendAllMissing(others LocalBranchNames) LocalBranchNames {
	missing := slice.FindAllMissing(self, others)
	slice.NaturalSort(missing)
	return append(self, missing...)
}

func (self LocalBranchNames) BranchNames() []BranchName {
	result := make([]BranchName, len(self))
	for l, localBranchName := range self {
		result[l] = localBranchName.BranchName()
	}
	return result
}

// Contains indicates whether this collection contains the given branch.
func (self LocalBranchNames) Contains(branch LocalBranchName) bool {
	return slices.Contains(self, branch)
}

// Hoist returns the given list with the given needles moved to the front,
// in the order they are given.
func (self LocalBranchNames) Hoist(branches ...LocalBranchName) LocalBranchNames {
	matches := make(LocalBranchNames, 0, len(branches))
	nonMatches := make(LocalBranchNames, 0, len(self)-len(branches))
	for _, branch := range self {
		if slices.Contains(branches, branch) {
			matches = append(matches, branch)
		} else {
			nonMatches = append(nonMatches, branch)
		}
	}
	return append(matches, nonMatches...)
}

// Join provides the names of all branches in this collection connected by the given separator.
func (self LocalBranchNames) Join(sep string) string {
	return strings.Join(self.Strings(), sep)
}

// Remove removes the given branch names from this collection.
func (self LocalBranchNames) Remove(toRemove ...LocalBranchName) LocalBranchNames {
	result := make(LocalBranchNames, 0, len(self))
	for _, branch := range self {
		if !slices.Contains(toRemove, branch) {
			result = append(result, branch)
		}
	}
	return result
}

func (self LocalBranchNames) String() string {
	if self == nil {
		return ""
	}
	return strings.Join(self.Strings(), " ")
}

// Strings provides the names of all branches in this collection as strings.
func (self LocalBranchNames) Strings() []string {
	result := make([]string, len(self))
	for b, branch := range self {
		result[b] = branch.String()
	}
	return result
}
