package common

import "github.com/git-town/git-town/v9/src/domain"

type SHAForBranchFunc func(domain.BranchName) (domain.SHA, error)
