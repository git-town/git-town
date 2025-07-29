package git

import (
	"fmt"
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
	// AUTO-RESOLVING CONFLICTS ENCOUNTERED WHILE SYNCING WITH THE TRACKING BRANCH
	//
	// When is it a proper rebase conflict?
	// When both the current branch and the tracking branch have different SHA compared to the end of the previous Git Town command that synced them.
	//
	// When is it a phantom rebase conflict?
	// When the tracking branch has the same SHA it had at the end of the previous Git Town command,
	// and the current branch has a different SHA --> the current branch was changed/rebased, we can force-push it.
	//
	// AUTO-RESOLVING CONFLICTS ENCOUNTERED WHILE SYNCING WITH THE PARENT BRANCH
	//
	// When is it a proper rebase conflict?
	// When both the current and the tracking branch have different SHA compared to the last time they were synced.
	//
	// When is it a phantom rebase conflict?
	// When the parent branch has the same SHA it had at the end of the previous Git Town command that synced it,
	// and the current branch has a different SHA compared to the end of the last Git Town command that synced it.
	// If there is a conflict now, we can keep the version of the current branch.
	//
	// We need to compare to the previous SHA when the branch was last synced.
	// If we just compare to the SHA at the end of the last command,
	// and the last command was "git town park", then it would wrongfully think the other branch has no changes.

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
		fmt.Println("1111111111111111111111111111111111111")
		fmt.Println("parent:", hasParentBlob, parentBlob)
		fmt.Println("current:", hasCurrentBlob, currentBlob)
		fmt.Println("root:", hasRootBlob, rootBlob)
		fmt.Println(fileConflict)
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
