package undobranches

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/undo/undodomain"
)

type LocalBranchRename undodomain.Change[gitdomain.LocalBranchName]
