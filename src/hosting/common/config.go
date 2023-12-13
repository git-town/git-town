package common

import "strings"

// Config contains data needed by all platform connectors.
type Config struct {
	// Hostname override
	Hostname string

	// the Organization within the hosting platform that owns the repo
	Organization string

	// repo name within the organization
	Repository string
}

func (self Config) HostnameWithStandardPort() string {
	index := strings.IndexRune(self.Hostname, ':')
	if index == -1 {
		return self.Hostname
	}
	return self.Hostname[:index]
}
