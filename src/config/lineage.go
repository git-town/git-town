package config

import "sort"

type BranchWithParent struct {
	Name   string
	Parent string
}

type Lineage []BranchWithParent

func (l Lineage) Lookup(name string) *BranchWithParent {
	for b, branchWithParent := range l {
		if branchWithParent.Name == name {
			return &l[b]
		}
	}
	return nil
}

// Ancestors provides all branches that are (great)(grand)parents of the branch with the given name.
func (l Lineage) Ancestors(branch string) Lineage {
	result := Lineage{}
	current := l.Lookup(branch)
	if current == nil {
		return result
	}
	for {
		parent := l.Lookup(current.Parent)
		if parent == nil {
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

// Children provides all branches that have the branch with the given name as their parent.
func (l Lineage) Children(branchName string) Lineage {
	result := Lineage{}
	for _, b := range l {
		if b.Parent == branchName {
			result = append(result, b)
		}
	}
	sort.Slice(result, func(a, b int) bool { return result[a].Name < result[b].Name })
	return result
}

// Contains indicates whether this Lineage contains a branch with the given name
func (l Lineage) Contains(branchName string) bool {
	for _, branch := range l {
		if branch.Name == branchName {
			return true
		}
	}
	return false
}

// HasParents returns whether or not the given branch has at least one parent.
func (l Lineage) HasParent(branch string) bool {
	for _, branchInfo := range l {
		if branchInfo.Parent == branch {
			return true
		}
	}
	return false
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (l Lineage) IsAncestor(ancestor, other string) bool {
	current := l.Lookup(other)
	if current == nil {
		return false
	}
	for {
		parent := l.Lookup(current.Parent)
		if parent == nil {
			return false
		}
		if parent.Name == ancestor {
			return true
		}
		current = parent
	}
}

// Roots provides the branches without parents, i.e. branches that start a branch lineage.
func (l Lineage) Roots() Lineage {
	roots := Lineage{}
	for _, branch := range l {
		if branch.Parent == "" && !roots.Contains(branch.Name) {
			roots = append(roots, branch)
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
		return l.IsAncestor(result[x].Parent, result[y].Parent)
	})
	return result
}
