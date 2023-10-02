package domain

import (
	"github.com/git-town/git-town/v9/src/slice"
)

// Remotes answers questions which Git remotes a repo has.
type Remotes []Remote

func NewRemotes(remotes ...string) Remotes {
	result := make(Remotes, len(remotes))
	for r, remote := range remotes {
		result[r] = NewRemote(remote)
	}
	return result
}

func (rs Remotes) HasOrigin() bool {
	return slice.Contains(rs, OriginRemote)
}

func (rs Remotes) HasUpstream() bool {
	return slice.Contains(rs, UpstreamRemote)
}
