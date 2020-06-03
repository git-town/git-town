package drivers

// UnsupportedHostingServiceError represents the error condition
// when no suitable hosting service is found in the respective Registry.
type UnsupportedHostingServiceError struct {
	registry *Registry
}

func (e UnsupportedHostingServiceError) Error() string {
	return `Unsupported hosting service

This command requires hosting on one of these services:
* Bitbucket
* GitHub
* GitLab
`
}
