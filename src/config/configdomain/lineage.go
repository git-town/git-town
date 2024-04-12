package configdomain

import (
	"sort"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/gohacks/slice"
	"golang.org/x/exp/maps"
)

// Lineage encapsulates all data and functionality around parent branches.
// branch --> its parent
// Lineage only contains branches that have ancestors.
type Lineage map[gitdomain.LocalBranchName]gitdomain.LocalBranchName

// Ancestors provides the names of all parent branches of the branch with the given name.
func (self Lineage) Ancestors(branch gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	current := branch
	result := gitdomain.LocalBranchNames{}
	for {
		parent, found := self[current]
		if !found {
			return result
		}
		result = append(gitdomain.LocalBranchNames{parent}, result...)
		current = parent
	}
}

// AncestorsWithoutRoot provides the names of all parent branches of the branch with the given name, excluding the root perennial branch.
func (self Lineage) AncestorsWithoutRoot(branch gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	current := branch
	result := gitdomain.LocalBranchNames{}
	for {
		parent, found := self[current]
		if !found {
			if len(result) == 0 {
				return result
			}
			return result[1:]
		}
		result.Prepend(parent)
		current = parent
	}
}

// BranchAndAncestors provides the full ancestry for the branch with the given name,
// including the branch.
func (self Lineage) BranchAndAncestors(branchName gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	return append(self.Ancestors(branchName), branchName)
}

// BranchLineageWithoutRoot provides all branches in the lineage of the given branch,
// from oldest to youngest, including the given branch.
func (self Lineage) BranchLineageWithoutRoot(branch gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	if self.Parent(branch).IsEmpty() {
		return gitdomain.LocalBranchNames{}
	}
	result := append(self.AncestorsWithoutRoot(branch), branch)
	return append(result, self.Descendants(branch)...)
}

// BranchNames provides the names of all branches in this Lineage, sorted alphabetically.
func (self Lineage) BranchNames() gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames(maps.Keys(self))
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
	self.OrderHierarchically(result)
	return result
}

// Children provides the names of all branches that have the given branch as their parent.
func (self Lineage) Children(branch gitdomain.LocalBranchName) gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
	for child, parent := range self {
		if parent == branch {
			result = append(result, child)
		}
	}
	result.Sort()
	return result
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

// HasParents returns whether or not the given branch has at least one parent.
func (self Lineage) HasParents(branch gitdomain.LocalBranchName) bool {
	for child := range self {
		if child == branch {
			return true
		}
	}
	return false
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (self Lineage) IsAncestor(ancestor, other gitdomain.LocalBranchName) bool {
	current := other
	for {
		parent, found := self[current]
		if !found {
			return false
		}
		if parent == ancestor {
			return true
		}
		current = parent
	}
}

// OrderHierarchically sorts the given branches in place so that ancestor branches come before their descendants
// and everything is sorted alphabetically.
func (self Lineage) OrderHierarchically(branches gitdomain.LocalBranchNames) {
	sort.Slice(branches, func(a, b int) bool {
		first := branches[a]
		second := branches[b]
		if first.IsEmpty() {
			return true
		}
		if second.IsEmpty() {
			return false
		}
		if self.IsAncestor(first, second) {
			return true
		}
		if self.IsAncestor(second, first) {
			return false
		}
		return first.String() < second.String()
	})
}

// Parent provides the name of the parent branch for the given branch or nil if the branch has no parent.
func (self Lineage) Parent(branch gitdomain.LocalBranchName) gitdomain.LocalBranchName {
	for child, parent := range self {
		if child == branch {
			return parent
		}
	}
	return gitdomain.EmptyLocalBranchName()
}

// RemoveBranch removes the given branch completely from this lineage.
func (self Lineage) RemoveBranch(branch gitdomain.LocalBranchName) {
	parent := self.Parent(branch)
	for _, childName := range self.Children(branch) {
		if parent.IsEmpty() {
			delete(self, childName)
		} else {
			self[childName] = parent
		}
	}
	delete(self, branch)
}

// Roots provides the branches with children and no parents.
func (self Lineage) Roots() gitdomain.LocalBranchNames {
	roots := gitdomain.LocalBranchNames{}
	for _, parent := range self {
		_, found := self[parent]
		if !found && !slice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	roots.Sort()
	return roots
}
