package configdomain

import (
	"cmp"
	"maps"
	"slices"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/mapstools"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Lineage encapsulates all data and functionality around parent branches.
// branch --> its parent
// Lineage only contains branches that have ancestors.
type Lineage struct {
	data LineageData
}

type LineageData map[gitdomain.LocalBranchName]gitdomain.LocalBranchName

func NewLineage() Lineage {
	return Lineage{
		data: make(LineageData),
	}
}

func NewLineageWith(data LineageData) Lineage {
	return Lineage{
		data: data,
	}
}

// Ancestors provides the names of all parent branches of the branch with the given name.
func (self Lineage) Ancestors(branch gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	current := branch
	result := gitdomain.LocalBranchNames{}
	for {
		parent, found := self.data[current]
		if !found {
			return result
		}
		result = append(gitdomain.LocalBranchNames{parent}, result...)
		current = parent
	}
}

// AncestorsWithoutRoot provides the names of all parent branches of the branch with the given name, excluding the root perennial branch.
func (self Lineage) AncestorsWithoutRoot(branch gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	ancestors := self.Ancestors(branch)
	if len(ancestors) == 0 {
		return ancestors
	}
	return ancestors[1:]
}

// BranchAndAncestors provides the full ancestry for the branch with the given name,
// including the branch.
func (self Lineage) BranchAndAncestors(branchName gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	return append(self.Ancestors(branchName), branchName)
}

// BranchAndAncestorsWithoutRoot provides the full ancestry for the branch with the given name,
// including the branch, but without the perennial branch at the root.
func (self Lineage) BranchAndAncestorsWithoutRoot(branchName gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	return append(self.AncestorsWithoutRoot(branchName), branchName)
}

// BranchLineageWithoutRoot provides all branches in the lineage of the given branch,
// from oldest to youngest, including the given branch.
func (self Lineage) BranchLineageWithoutRoot(branch gitdomain.LocalBranchName, perennialBranches gitdomain.LocalBranchNames) gitdomain.LocalBranchNames {
	if self.Parent(branch).IsNone() {
		result := gitdomain.LocalBranchNames{}
		if !perennialBranches.Contains(branch) {
			result = append(result, branch)
		}
		return append(result, self.Descendants(branch)...)
	}
	return append(append(self.AncestorsWithoutRoot(branch), branch), self.Descendants(branch)...)
}

// BranchNames provides the names of all branches in this Lineage, sorted alphabetically.
func (self Lineage) BranchNames() gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames(slices.Collect(maps.Keys(self.data)))
	result.Sort()
	return result
}

// BranchesAndAncestors provides the full lineage for the branches with the given names,
// including the branches themselves.
func (self Lineage) BranchesAndAncestors(branchNames gitdomain.LocalBranchNames) gitdomain.LocalBranchNames {
	result := branchNames
	for _, branchName := range branchNames {
		result = slice.AppendAllMissing(result, self.Ancestors(branchName)...)
	}
	return self.OrderHierarchically(result)
}

// provides all branches for which the parent is known
func (self Lineage) BranchesWithParents() gitdomain.LocalBranchNames {
	var result gitdomain.LocalBranchNames = slices.Collect(maps.Keys(self.data))
	result.Sort()
	return result
}

// Children provides the names of all branches that have the given branch as their parent.
func (self Lineage) Children(branch gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
	for child, parent := range self.data {
		if parent == branch {
			result = append(result, child)
		}
	}
	return slice.NaturalSort(result)
}

// Descendants provides all branches that depend on the given branch in its lineage.
func (self Lineage) Descendants(branch gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
	for _, child := range self.Children(branch) {
		result = append(result, child)
		result = append(result, self.Descendants(child)...)
	}
	return result
}

func (self Lineage) Entries() []LineageEntry {
	result := make([]LineageEntry, 0, self.Len())
	for branch, parent := range self.data {
		result = append(result, LineageEntry{
			Child:  branch,
			Parent: parent,
		})
	}
	slices.SortFunc(result, func(a, b LineageEntry) int {
		return cmp.Compare(a.Child, b.Child)
	})
	return result
}

func (self Lineage) HasDescendents(branch gitdomain.LocalBranchName) bool {
	for _, parent := range self.data {
		if parent == branch {
			return true
		}
	}
	return false
}

// HasParents returns whether or not the given branch has at least one parent.
func (self Lineage) HasParents(branch gitdomain.LocalBranchName) bool {
	for child := range self.data {
		if child == branch {
			return true
		}
	}
	return false
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (self Lineage) IsAncestor(ancestor, other gitdomain.LocalBranchName) bool {
	ancestors := self.Ancestors(other)
	return ancestors.Contains(ancestor)
}

func (self Lineage) IsEmpty() bool {
	return self.data == nil || self.Len() == 0
}

// provides the branch from candidates that is the youngest ancestor of the given branch
func (self Lineage) LatestAncestor(branch gitdomain.LocalBranchName, candidates gitdomain.LocalBranchNames) Option[gitdomain.LocalBranchName] {
	for {
		if candidates.Contains(branch) {
			return Some(branch)
		}
		if parent, hasParent := self.Parent(branch).Get(); hasParent {
			branch = parent
		} else {
			return None[gitdomain.LocalBranchName]()
		}
	}
}

func (self Lineage) Len() int {
	return len(self.data)
}

// provides a new Lineage that consists of entries from both this and the given Lineage
func (self Lineage) Merge(other Lineage) Lineage {
	return Lineage{
		data: mapstools.Merge(self.data, other.data),
	}
}

// OrderHierarchically provides the given branches sorted so that ancestor branches come before their descendants.
func (self Lineage) OrderHierarchically(branches gitdomain.LocalBranchNames) gitdomain.LocalBranchNames {
	result := make(gitdomain.LocalBranchNames, 0, len(self.data))
	for _, root := range self.Roots() {
		self.addChildrenHierarchically(&result, root, branches)
	}
	result = result.AppendAllMissing(branches...)
	return result
}

// Parent provides the name of the parent branch for the given branch or nil if the branch has no parent.
func (self Lineage) Parent(branch gitdomain.LocalBranchName) Option[gitdomain.LocalBranchName] {
	for child, parent := range self.data {
		if child == branch {
			return Some(parent)
		}
	}
	return None[gitdomain.LocalBranchName]()
}

// RemoveBranch removes the given branch completely from this lineage.
func (self Lineage) RemoveBranch(branch gitdomain.LocalBranchName) Lineage {
	delete(self.data, branch)
	return self
}

// Root provides the oldest ancestor of the given branch.
func (self Lineage) Root(branch gitdomain.LocalBranchName) gitdomain.LocalBranchName {
	current := branch
	for {
		parent, found := self.data[current]
		if !found {
			return current
		}
		current = parent
	}
}

// Roots provides the branches with children and no parents.
func (self Lineage) Roots() gitdomain.LocalBranchNames {
	roots := gitdomain.LocalBranchNames{}
	for _, parent := range self.data {
		_, found := self.data[parent]
		if !found && !slices.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	roots.Sort()
	return roots
}

func (self Lineage) Set(branch, parent gitdomain.LocalBranchName) Lineage {
	self = self.initializeIfNeeded()
	self.data[branch] = parent
	return self
}

func (self Lineage) addChildrenHierarchically(result *gitdomain.LocalBranchNames, currentBranch gitdomain.LocalBranchName, allBranches gitdomain.LocalBranchNames) {
	if allBranches.Contains(currentBranch) {
		*result = append(*result, currentBranch)
	}
	for _, child := range self.Children(currentBranch) {
		self.addChildrenHierarchically(result, child, allBranches)
	}
}

func (self Lineage) initializeIfNeeded() Lineage {
	if self.data == nil {
		self.data = make(map[gitdomain.LocalBranchName]gitdomain.LocalBranchName)
	}
	return self
}
