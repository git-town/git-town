package git

import (
	"reflect"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

func DetectPhantomRebaseConflicts(conflictInfos []FileConflictFullInfo, parentBranch gitdomain.BranchName, rootBranch gitdomain.LocalBranchName) []PhantomConflict {
	// How to detect phantom merge conflicts:
	//
	// Situations we want to cover:
	//
	// 1. Feature branch changed (commit added, amended, or rebased) and root not changed --> Keep the feature branch change
	// 2. Root branch branch was changed and feature branch not changed --> don't auto-resolve
	// 3. Both happen at the same time: root and feature branch changed --> don't auto-resolve
	//
	// Determine upfront:
	// Feature branch is changed (current SHA is different from SHA at end of last run)
	// Root branch is changed
	//
	// One side is the old feature branch: auto-resolve to the other side?
	// One side is the old root branch: auto-resolve to the other side?
	//
	// O
	//
	// if parentBranch == rootBranch.BranchName() {
	// 	// branches whose parent is the root branch cannot have phantom merge conflicts
	// 	return []PhantomMergeConflict{}
	// }
	result := []PhantomConflict{}
	for _, conflictInfo := range conflictInfos {
		// TODO: inspect the conflictInfo
		initialParentInfo, hasInitialParentInfo := conflictInfo.Parent.Get()
		currentInfo, hasCurrentInfo := conflictInfo.Current.Get()
		if !hasInitialParentInfo || !hasCurrentInfo || currentInfo.Permission != initialParentInfo.Permission {
			continue
		}
		if reflect.DeepEqual(conflictInfo.Root, conflictInfo.Parent) {
			// root and parent have the exact same version of the file --> this is a phantom rebase conflict
			result = append(result, PhantomConflict{
				FilePath: currentInfo.FilePath,
			})
		}
	}
	return result
}
