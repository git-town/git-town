package undobranches

import (
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/undo/undodomain"
)

type LocalBranchRename undodomain.Change[gitdomain.LocalBranchName]
