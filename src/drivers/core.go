package drivers

import (
	"errors"

	"github.com/git-town/git-town/src/git"
)

// Get returns the code hosting driver to use based on the git config.
func Get(config git.Configuration) (CodeHostingDriver, error) {
	if driver := tryCreateBitBucket(config); driver != nil {
		return driver, nil
	}
	if driver := tryCreateGitHub(config); driver != nil {
		return driver, nil
	}
	if driver := tryCreateGitLab(config); driver != nil {
		return driver
	}
	if driver := tryCreateGitea(config); driver != nil {
		return driver
	}
	return nil, errors.New(`unsupported hosting service

This command requires hosting on one of these services:
* Gitea
* GitHub
* GitLab
* Bitbucket
`)
}
