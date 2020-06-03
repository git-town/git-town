package drivers

import "github.com/git-town/git-town/src/git"

// DriverManager is the public API of the drivers package.
type DriverManager struct {
	bitBucket bitbucketCodeHostingDriver
	gitHub    githubCodeHostingDriver
	gitLab    gitlabCodeHostingDriver
}

// ActiveDriver returns the code hosting driver to use based on the git config.
func (dm *DriverManager) ActiveDriver(config git.Configuration) CodeHostingDriver {
	if driver := tryCreateBitBucket(config); driver != nil {
		return driver
	}
	if driver := tryCreateGitHub(config); driver != nil {
		return driver
	}
	if driver := dm.gitLab.TryCreate(config); driver != nil {
		return driver
	}
	return nil
}
