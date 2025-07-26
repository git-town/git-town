package git

import (
	"reflect"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

func DetectPhantomRebaseConflicts(conflictInfos []FileConflictFullInfo, parentBranch gitdomain.BranchName, rootBranch gitdomain.LocalBranchName) []PhantomConflict {
	// How to detect phantom merge conflicts:
	//
	// One side is the root branch from the last run: this was resolved before, automatically resolve using the other side
	//   - can this even happen?
	//     One side is the current root branch, and its different from the root branch of the last run --> don't auto-resolve
	// One side is the feature branch from the last run: this was
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
