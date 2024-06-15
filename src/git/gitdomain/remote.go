package gitdomain

import "github.com/git-town/git-town/v14/src/gohacks"

// Remote represents a Git remote.
type Remote gohacks.NonEmptyString

func NewRemote(id string) Remote {
	for _, remote := range AllRemotes {
		if id == remote.String() {
			return remote
		}
	}
	return RemoteOther
}

// Implementation of the fmt.Stringer interface.
func (self Remote) String() string {
	return string(self)
}

const (
	RemoteNone     = Remote("")
	RemoteOrigin   = Remote("origin")
	RemoteOther    = Remote("other")
	RemoteUpstream = Remote("upstream")
)

var AllRemotes = []Remote{ //nolint:gochecknoglobals
	RemoteNone,
	RemoteOrigin,
	RemoteOther,
	RemoteUpstream,
}
