// Package cache provides infrastructure to cache things in memory.
package cache

import "github.com/git-town/git-town/v9/src/domain"

// Bool is a cache for bool variables.
type Bool = Cache[bool]

// LocalBranch is a cache for domain.LocalBranchName variables.
type LocalBranch = Cache[domain.LocalBranchName]

// RemoteBranch is a cache for domain.RemoteBranchName variables.
type RemoteBranch = Cache[domain.RemoteBranchName]

// Remotes is a cache for domain.Remotes variables.
type Remotes = Cache[domain.Remotes]

// String is a cache for string variables.
type String = Cache[string]

// Strings is a cache for string variables.
type Strings = Cache[[]string]
