package drivers

import "github.com/git-town/git-town/src/git"

// Get returns the code hosting driver to use based on the git config.
func Get(config git.Configuration) CodeHostingDriver {
	if driver := tryCreateBitBucket(config); driver != nil {
		return driver
	}
	if driver := tryCreateGitHub(config); driver != nil {
		return driver
	}
	if driver := tryCreateGitLab(config); driver != nil {
		return driver
	}
	return nil
}
