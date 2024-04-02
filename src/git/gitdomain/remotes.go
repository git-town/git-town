package gitdomain

import (
	"github.com/git-town/git-town/v13/src/gohacks/slice"
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

func (self Remotes) HasOrigin() bool {
	return slice.Contains(self, RemoteOrigin)
}

func (self Remotes) HasUpstream() bool {
	return slice.Contains(self, RemoteUpstream)
}
