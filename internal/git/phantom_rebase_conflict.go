package git

import (
	"fmt"
	"reflect"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

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
