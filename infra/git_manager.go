package infra

import (
	"path"

	"github.com/pkg/errors"
)

/*
GitManager manages the Git setup for the entire test suite.
In particular, it creates the Git setup for individual feature specs.
*/
type GitManager struct {

	// path of the folder that this class operates in
	dir string

	// the memoized environment
	memoized *GitEnvironment
}

// NewGitManager creates a new GitManager instance.
func NewGitManager(baseDir string) *GitManager {
	return &GitManager{dir: baseDir}
}

// CreateMemoizedEnvironment creates the Git environment cache
// that makes cloning new GitEnvironment instances faster.
func (gm *GitManager) CreateMemoizedEnvironment() error {
	var err error
	gm.memoized, err = NewGitEnvironment(path.Join(gm.dir, "memoized"))
	if err != nil {
		return errors.Wrapf(err, "cannot create memoized environment")
	}
	err = gm.memoized.Populate()
	if err != nil {
		return errors.Wrapf(err, "cannot populate memoized environment")
	}
	return nil
}

// CreateScenarioEnvironment creates a new GitEnvironment for the scenario with the given name
func (gm GitManager) CreateScenarioEnvironment(scenarioName string) (*GitEnvironment, error) {
	envPath := path.Join(gm.dir, scenarioName)
	return CloneGitEnvironment(gm.memoized, envPath)
}
