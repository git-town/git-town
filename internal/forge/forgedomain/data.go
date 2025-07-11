package forgedomain

import "strings"

// Data contains data needed by all forge connectors.
type Data struct {
	// Hostname override
	Hostname string

	// the Organization within the forge that owns the repo
	Organization string

	// repo name within the organization
	Repository string
}

func (self Data) HostnameWithStandardPort() string {
	index := strings.IndexRune(self.Hostname, ':')
	if index == -1 {
		return self.Hostname
	}
	return self.Hostname[:index]
}
