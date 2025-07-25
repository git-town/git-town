package git

import (
	"reflect"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

func DetectPhantomRebaseConflicts(conflictInfos []FileConflictFullInfo, parentBranch gitdomain.BranchName, rootBranch gitdomain.LocalBranchName) []PhantomMergeConflict {
	// TODO: inspect the conflictInfos here

	if parentBranch == rootBranch.BranchName() {
		// branches whose parent is the root branch cannot have phantom merge conflicts
		return []PhantomMergeConflict{}
	}
	result := []PhantomMergeConflict{}
	for _, conflictInfo := range conflictInfos {
		initialParentInfo, hasInitialParentInfo := conflictInfo.Parent.Get()
		currentInfo, hasCurrentInfo := conflictInfo.Current.Get()
		if !hasInitialParentInfo || !hasCurrentInfo || currentInfo.Permission != initialParentInfo.Permission {
			continue
		}
		if reflect.DeepEqual(conflictInfo.Root, conflictInfo.Parent) {
			// root and parent have the exact same version of the file --> this is a phantom merge conflict
			result = append(result, PhantomMergeConflict{
				FilePath: currentInfo.FilePath,
			})
		}
	}
	return result
}
