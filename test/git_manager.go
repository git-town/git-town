package test

import (
	"math/rand"
	"path"
	"strconv"

	"github.com/pkg/errors"
)

/*
GitManager manages the Git setup for the entire test suite.
In particular, it creates the Git setup for individual feature specs (GitEnvironment).
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
func (manager *GitManager) CreateMemoizedEnvironment() error {
	var err error
	manager.memoized, err = NewStandardGitEnvironment(path.Join(manager.dir, "memoized"))
	if err != nil {
		return errors.Wrapf(err, "cannot create memoized environment")
	}
	return nil
}

// CreateScenarioEnvironment creates a new GitEnvironment for the scenario with the given name
func (manager *GitManager) CreateScenarioEnvironment(scenarioName string) (*GitEnvironment, error) {
	envPath := path.Join(manager.dir, strconv.Itoa(rand.Intn(9999))+scenarioName)
	return CloneGitEnvironment(manager.memoized, envPath)
}
