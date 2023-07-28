package config

import "github.com/git-town/git-town/v9/src/stringslice"

// Remotes contains the name of the remotes that are set up in the Git repo.
type Remotes []string

func (r Remotes) HasOrigin() bool {
	return stringslice.Contains(r, OriginRemote)
}

func (r Remotes) HasUpstream() bool {
	return stringslice.Contains(r, UpstreamRemote)
}

// OriginRemote contains the name of the "origin" remote.
const OriginRemote = "origin"

// UpstreamRemote contains the name of the "upstream" remote.
const UpstreamRemote = "upstream"
