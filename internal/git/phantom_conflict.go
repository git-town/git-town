package git

import "github.com/git-town/git-town/v21/internal/git/gitdomain"

// describes a file within an unresolved merge conflict that experiences a phantom merge conflict
type PhantomConflict struct {
	FilePath   string
	Resolution gitdomain.ConflictResolution
}
