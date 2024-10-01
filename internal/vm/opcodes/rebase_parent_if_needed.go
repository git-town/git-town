package opcodes

import "github.com/git-town/git-town/v16/internal/git/gitdomain"

type RebaseParentIfNeeded struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}
