package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// HostedRepoInfo provides information about a repository hosted at a forge.
type HostedRepoInfo struct {
	// hostname of the server that hosts the repo
	Hostname string

	// the Organization within the forge that owns the repo
	Organization string

	// name of the repo on the server
	Repository string

	Supergroup Option[string]
}

func (self HostedRepoInfo) HostnameWithStandardPort() string {
	index := strings.IndexRune(self.Hostname, ':')
	if index == -1 {
		return self.Hostname
	}
	return self.Hostname[:index]
}
