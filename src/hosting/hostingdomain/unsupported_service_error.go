package hostingdomain

import "errors"

// UnsupportedServiceError communicates that the origin remote runs an unknown code hosting platform.
func UnsupportedServiceError() error {
	return errors.New(`unsupported hosting platform

This command requires hosting on one of these services:
* Bitbucket
* GitHub
* GitLab
* Gitea`)
}
