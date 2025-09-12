package forgedomain

import "strings"

// HostedRepoInfo provides information about a repository hosted at a forge.
type HostedRepoInfo struct {
	// hostname of the server that hosts the repo
	Hostname string

	// the Organization within the forge that owns the repo
	Organization string

	// name of the repo on the server
	Repository string
}

func (self HostedRepoInfo) HostnameWithStandardPort() string {
	index := strings.IndexRune(self.Hostname, ':')
	if index == -1 {
		return self.Hostname
	}
	return self.Hostname[:index]
}
