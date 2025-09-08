package forgedomain

import "errors"

// UnsupportedServiceError communicates that the origin remote runs an unknown forge type.
func UnsupportedServiceError() error {
	return errors.New(`unsupported forge type

This command requires hosting on one of these services:
* Bitbucket
* Bitbucket Data Center
* Forgejo
* GitHub
* GitLab
* Gitea`)
}
