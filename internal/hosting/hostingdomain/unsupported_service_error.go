package hostingdomain

import "errors"

// UnsupportedServiceError communicates that the origin remote runs an unknown forge.
func UnsupportedServiceError() error {
	return errors.New(`unsupported hosting platform

This command requires hosting on one of these services:
* Bitbucket
* Bitbucket Data Center
* GitHub
* GitLab
* Gitea`)
}
