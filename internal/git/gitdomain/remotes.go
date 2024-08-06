package gitdomain

import "slices"

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
	return slices.Contains(self, RemoteOrigin)
}

func (self Remotes) HasUpstream() bool {
	return slices.Contains(self, RemoteUpstream)
}
