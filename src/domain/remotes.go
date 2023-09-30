package domain

import "github.com/git-town/git-town/v9/src/comparables"

// Remotes answers questions which Git remotes a repo has.
type Remotes []Remote

func NewRemotes(remotes ...string) Remotes {
	result := make(Remotes, len(remotes))
	for r, remote := range remotes {
		result[r] = NewRemote(remote)
	}
	return result
}

func (r Remotes) HasOrigin() bool {
	return comparables.Contains(r, OriginRemote)
}

func (r Remotes) HasUpstream() bool {
	return comparables.Contains(r, UpstreamRemote)
}
