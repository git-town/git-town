package datatable

import (
	"sort"

	prodgit "github.com/git-town/git-town/v8/src/git"
)

// BranchHierarchyTable provides the currently configured branch hierarchy information as a DataTable.
func BranchHierarchyTable(config *prodgit.RepoConfig) DataTable {
	result := DataTable{}
	config.Reload()
	parentBranchMap := config.ParentBranchMap()
	result.AddRow("BRANCH", "PARENT")
	childBranches := make([]string, 0, len(parentBranchMap))
	for child := range parentBranchMap {
		childBranches = append(childBranches, child)
	}
	sort.Strings(childBranches)
	for _, child := range childBranches {
		result.AddRow(child, parentBranchMap[child])
	}
	return result
}
