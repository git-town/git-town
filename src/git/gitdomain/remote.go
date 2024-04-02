package gitdomain

import "fmt"

// Remote represents a Git remote.
type Remote string

func NewRemote(id string) Remote {
	for _, remote := range AllRemotes {
		if id == remote.String() {
			return remote
		}
	}
	panic(fmt.Sprintf("unknown remote: %q", id))
}

func (self Remote) IsEmpty() bool {
	return self == ""
}

// Implementation of the fmt.Stringer interface.
func (self Remote) String() string {
	return string(self)
}

const (
	RemoteNone     = Remote("")
	RemoteOrigin   = Remote("origin")
	RemoteUpstream = Remote("upstream")
)

var AllRemotes = []Remote{ //nolint:gochecknoglobals
	RemoteNone,
	RemoteOrigin,
	RemoteUpstream,
}
