package git

import (
	"reflect"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// PhantomConflict describes a file within an unresolved merge conflict that experiences a phantom merge conflict.
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
