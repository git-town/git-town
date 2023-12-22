package hostingdomain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// TODO: remove and replace with a direct link to git.BackendCommands.
type SHAForBranchFunc func(gitdomain.BranchName) (gitdomain.SHA, error)
