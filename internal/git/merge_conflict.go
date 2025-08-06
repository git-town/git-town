package git

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Everything Git Town needs to know about a file merge conflict to determine whether this is a phantom merge conflict.
type MergeConflict struct {
	Current Option[Blob] // info about the file on the current branch
	Parent  Option[Blob] // info about the file on the original parent
	Root    Option[Blob] // info about the file on the root branch
}

func (mergeConflict MergeConflict) Debug(querier subshelldomain.Querier) {
	current, hasCurrent := mergeConflict.Current.Get()
	parent, hasParent := mergeConflict.Parent.Get()
	root, hasRoot := mergeConflict.Root.Get()
	fmt.Print("ROOT: ")
	if hasRoot {
		root.Debug(querier)
	} else {
		fmt.Println("(none)")
	}
	fmt.Print("PARENT CHANGE: ")
	if hasParent {
		parent.Debug(querier)
	} else {
		fmt.Println("(none)")
	}
	fmt.Print("CURRENT CHANGE: ")
	if hasCurrent {
		current.Debug(querier)
	} else {
		fmt.Println("(none)")
	}
}
