// Package cache provides infrastructure to cache things in memory.
package cache

import "github.com/git-town/git-town/v17/internal/git/gitdomain"

// LocalBranch is a cache for gitdomain.LocalBranchName variables.
type LocalBranchWithPrevious = WithPrevious[gitdomain.LocalBranchName]
