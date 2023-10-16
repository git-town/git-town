package common

// Config contains data needed by all platform connectors.
type Config struct {
	// bearer token to authenticate with the API
	APIToken string

	// Hostname override
	Hostname string

	// the Organization within the hosting platform that owns the repo
	Organization string

	// repo name within the organization
	Repository string
}
