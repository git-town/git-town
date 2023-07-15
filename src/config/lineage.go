package config

import (
	"sort"
)

type BranchWithParent struct {
	Name   string
	Parent string
}

type Lineage []BranchWithParent

// Ancestors provides all branches that are (great)(grand)parents of the branch with the given name.
func (l Lineage) Ancestors(branchName string) Lineage {
	result := Lineage{}
	current := l.Lookup(branchName)
	for {
		if current.Parent == "" {
			return result
		}
		parent := l.Lookup(current.Parent)
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

// IsAncestor indicates whether the branch with the given ancestor name is indeed an ancestor
// of the other branch with the given name.
func (l Lineage) IsAncestor(ancestorName, branchName string) bool {
	current := l.Lookup(branchName)
	for {
		if current == nil || current.Parent == "" {
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

// Lookup provides a copy of the lineage entry with the given name.
func (l Lineage) Lookup(name string) *BranchWithParent {
	for b, branchWithParent := range l {
		if branchWithParent.Name == name {
			return &l[b]
		}
	}
	return nil
}

func (l Lineage) Parent(branchName string) string {
	parentBranch := l.Lookup(branchName)
	if parentBranch == nil {
		return ""
	}
	return parentBranch.Parent
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
		return l.IsAncestor(result[x].Name, result[y].Name)
		// fmt.Printf("%s %s %v", result[x].Name)
	})
	return result
}
