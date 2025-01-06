package gitdomain

import "slices"

// Remotes answers questions which Git remotes a repo has.
type Remotes []Remote

func NewRemotes(remotes ...string) Remotes {
	result := make(Remotes, 0, len(remotes))
	for _, remoteName := range remotes {
		if remote, hasRemote := NewRemote(remoteName).Get(); hasRemote {
			result = append(result, remote)
		}
	}
	return result
}

func (self Remotes) Contains(remote Remote) bool {
	return slices.Contains(self, remote)
}

func (self Remotes) HasDev(devRemote Remote) bool {
	return self.Contains(devRemote)
}

func (self Remotes) HasUpstream() bool {
	return slices.Contains(self, RemoteUpstream)
}
