package undobranches

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/undo/undodomain"
)

type LocalBranchRename undodomain.Change[gitdomain.LocalBranchName]
