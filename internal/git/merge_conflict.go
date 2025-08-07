package git

import . "github.com/git-town/git-town/v21/pkg/prelude"

// Everything Git Town needs to know about a merge conflict to determine whether this is a phantom merge conflict.
type MergeConflict struct {
	Current Option[Blob] // info about the file on the current branch
	Parent  Option[Blob] // info about the file on the original parent
	Root    Option[Blob] // info about the file on the root branch
}
