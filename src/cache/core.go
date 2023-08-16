// Package cache provides infrastructure to cache things in memory.
package cache

import "github.com/git-town/git-town/v9/src/domain"

// Bool is a cache for bool variables.
type Bool = Cache[bool]

type LocalBranch = Cache[domain.LocalBranchName]

type RemoteBranch = Cache[domain.RemoteBranchName]

// String is a cache for string variables.
type String = Cache[string]

// Strings is a cache for string variables.
type Strings = Cache[[]string]
