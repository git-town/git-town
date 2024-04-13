// Package cache provides infrastructure to cache things in memory.
package cache

import "github.com/git-town/git-town/v14/src/git/gitdomain"

// Bool is a cache for bool variables.
type Bool = Cache[bool]

// LocalBranch is a cache for gitdomain.LocalBranchName variables.
type LocalBranchWithPrevious = WithPrevious[gitdomain.LocalBranchName]

// RemoteBranch is a cache for gitdomain.RemoteBranchName variables.
type RemoteBranch = Cache[gitdomain.RemoteBranchName]

// Remotes is a cache for domain.Remotes variables.
type Remotes = Cache[gitdomain.Remotes]

// String is a cache for string variables.
type String = Cache[string]

// Strings is a cache for string variables.
type Strings = Cache[[]string]
