package gitdomain

import (
	"slices"
	"sort"
	"strings"

	"github.com/git-town/git-town/v15/internal/gohacks/slice"
	. "github.com/git-town/git-town/v15/pkg/prelude"
)

type LocalBranchNames []LocalBranchName

func NewLocalBranchNames(names ...string) LocalBranchNames {
	result := make(LocalBranchNames, len(names))
	for n, name := range names {
		result[n] = NewLocalBranchName(name)
	}
	return result
}

// ParseLocalBranchNamesRef constructs a LocalBranchNames instance
// containing the branches listed in the given space-separated string.
func ParseLocalBranchNames(names string) LocalBranchNames {
	var branchNames []string
	if names != "" {
		branchNames = strings.Split(names, " ")
	}
	return NewLocalBranchNames(branchNames...)
}

// AppendAllMissing provides a LocalBranchNames list consisting of the sum of this and elements of other list that aren't in this list.
func (self LocalBranchNames) AppendAllMissing(others ...LocalBranchName) LocalBranchNames {
	return slice.AppendAllMissing(self, others...)
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

// Index provides the index of the element with the given text.
func (self LocalBranchNames) Index(needle LocalBranchName) Option[int] {
	for e, entry := range self {
		if entry == needle {
			return Some(e)
		}
	}
	return None[int]()
}

// IndexWithTextOr provides the index of the element with the given text
// or the given default index if the element isn't in this collection.
func (self LocalBranchNames) IndexOr(needle LocalBranchName, defaultIndex int) int {
	if index, found := self.Index(needle).Get(); found {
		return index
	}
	return defaultIndex
}

// Join provides the names of all branches in this collection connected by the given separator.
func (self LocalBranchNames) Join(sep string) string {
	return strings.Join(self.Strings(), sep)
}

func (self *LocalBranchNames) Prepend(branch LocalBranchName) {
	*self = append(LocalBranchNames{branch}, *self...)
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

// RemoveWorktreeMarkers removes the workspace markers from the branch names in this list.
func (self LocalBranchNames) RemoveWorktreeMarkers() LocalBranchNames {
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
