package common

import "github.com/git-town/git-town/v11/src/domain"

// TODO: remove and replace with a direct link to git.BackendCommands.
type SHAForBranchFunc func(domain.BranchName) (domain.SHA, error)
