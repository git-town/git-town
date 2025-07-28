package git

import (
	"reflect"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

func DetectPhantomRebaseConflicts(fileConflicts []FileConflictFullInfo, parentBranch gitdomain.BranchName, rootBranch gitdomain.LocalBranchName) []PhantomConflict {
	// OVERALL SYNC STRATEGY
	//
	// Sync tracking branch:
	//
	// Feature branch changed && tracking not changed --> overwrite tracking
	//
	// Feature branch changed && tracking changed --> sync normally
	//
	// Feature branch not changed && tracking changed --> pull
	//
	// Feature branch not changed && tracking not changed --> do nothing
	//
	// Sync parent branch:
	//
	// Feature branch changed && parent not changed --> do nothing
	//
	// Feature branch changed && parent changed --> sync normally
	//
	// Feature branch not changed && parent changed --> pull
	//
	// Feature branch not changed && parent not changed --> do nothing
	//
	// AUTO-RESOLVING CONFLICTS ENCOUNTERED WHILE SYNCING
	//
	//
	// O
	//
	// if parentBranch == rootBranch.BranchName() {
	// 	// branches whose parent is the root branch cannot have phantom merge conflicts
	// 	return []PhantomMergeConflict{}
	// }
	result := []PhantomConflict{}
	for _, fileConflict := range fileConflicts {
		// TODO: inspect the conflictInfo
		//
		// One side is the old feature branch: auto-resolve to the other side?
		// One side is the old root branch: auto-resolve to the other side?
		parentBlob, hasParentBlob := fileConflict.Parent.Get()
		currentBlob, hasCurrentBlob := fileConflict.Current.Get()
		rootBlob, hasRootBlob := fileConflict.Root.Get()
		if !hasParentBlob || !hasCurrentBlob || currentBlob.Permission != parentBlob.Permission {
			continue
		}
		if reflect.DeepEqual(fileConflict.Root, fileConflict.Parent) {
			// root and parent have the exact same version of the file --> this is a phantom rebase conflict
			result = append(result, PhantomConflict{
				FilePath: currentBlob.FilePath,
			})
		}
	}
	return result
}
