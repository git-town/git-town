package git

import (
	"fmt"
	"reflect"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// describes a file within an unresolved merge conflict that experiences a phantom merge conflict
type PhantomConflict struct {
	FilePath   string
	Resolution gitdomain.ConflictResolution
}

func DetectPhantomMergeConflicts(conflictInfos []MergeConflict, parentBranchOpt Option[gitdomain.LocalBranchName], rootBranch gitdomain.LocalBranchName) []PhantomConflict {
	parentBranch, hasParentBranch := parentBranchOpt.Get()
	if !hasParentBranch || parentBranch == rootBranch {
		// branches that don't have a parent or whose parent is the root branch cannot have phantom merge conflicts
		return []PhantomConflict{}
	}
	result := []PhantomConflict{}
	for _, conflictInfo := range conflictInfos {
		initialParentInfo, hasInitialParentInfo := conflictInfo.Parent.Get()
		currentInfo, hasCurrentInfo := conflictInfo.Current.Get()
		if !hasInitialParentInfo || !hasCurrentInfo || currentInfo.Permission != initialParentInfo.Permission {
			continue
		}
		if reflect.DeepEqual(conflictInfo.Root, conflictInfo.Parent) {
			// root and parent have the exact same version of the file --> this is a phantom merge conflict
			result = append(result, PhantomConflict{
				FilePath:   currentInfo.FilePath,
				Resolution: gitdomain.ConflictResolutionOurs,
			})
		}
	}
	return result
}

func DetectPhantomRebaseConflicts(fileConflicts MergeConflicts, parentBranch gitdomain.BranchName, rootBranch gitdomain.LocalBranchName) []PhantomConflict {
	// if parentBranch == rootBranch.BranchName() {
	// 	// branches whose parent is the root branch cannot have phantom merge conflicts
	// 	return []PhantomMergeConflict{}
	// }
	result := []PhantomConflict{}
	for _, fileConflict := range fileConflicts {
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
