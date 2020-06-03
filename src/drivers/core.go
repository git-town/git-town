package drivers

import (
	"errors"

	"github.com/git-town/git-town/src/git"
)

// Core provides the public API for the drivers subsystem.

// GetActiveDriver returns the code hosting driver to use based on the git config.
func GetActiveDriver() (CodeHostingDriver, error) {
	driver := TryUseGithub(git.Config())
	if driver != nil {
		return driver, nil
	}
	driver = TryUseBitbucket(git.Config())
	if driver != nil {
		return driver, nil
	}
	driver = TryUseGitlab(git.Config())
	if driver != nil {
		return driver, nil
	}
	return nil, errors.New(`unsupported hosting service

This command requires hosting on one of these services:
* Bitbucket
* GitHub
* GitLab`)
}
