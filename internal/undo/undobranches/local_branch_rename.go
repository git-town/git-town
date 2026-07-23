package undobranches

import (
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/undo/undodomain"
)

type LocalBranchRename undodomain.Change[gitdomain.LocalBranchName]
