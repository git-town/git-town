package config

import (
	"github.com/git-town/git-town/v9/src/genericslice"
)

// Remotes answers questions which Git remotes a repo has.
type Remotes []string

func (r Remotes) HasOrigin() bool {
	return genericslice.Contains(r, OriginRemote)
}

func (r Remotes) HasUpstream() bool {
	return genericslice.Contains(r, UpstreamRemote)
}

// OriginRemote contains the name of the "origin" remote.
const OriginRemote = "origin"

// UpstreamRemote contains the name of the "upstream" remote.
const UpstreamRemote = "upstream"
