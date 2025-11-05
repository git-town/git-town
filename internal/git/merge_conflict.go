package git

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// MergeConflict contains everything Git Town needs to know about a merge conflict to determine whether this is a phantom merge conflict.
type MergeConflict struct {
	Current Option[Blob] // info about the file on the current branch
	Parent  Option[Blob] // info about the file on the original parent
	Root    Option[Blob] // info about the file on the root branch
}

func (self MergeConflict) Debug(querier subshelldomain.Querier) {
	current, hasCurrent := self.Current.Get()
	parent, hasParent := self.Parent.Get()
	root, hasRoot := self.Root.Get()
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
