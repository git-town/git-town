package gitdomain

import (
	"strings"

	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Remote represents a Git remote.
type Remote string

func NewRemote(id string) Option[Remote] {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return None[Remote]()
	}
	return Some(Remote(id))
}

// Implementation of the fmt.Stringer interface.
func (self Remote) String() string {
	return string(self)
}

const (
	RemoteOrigin   = Remote("origin")
	RemoteUpstream = Remote("upstream")
)
