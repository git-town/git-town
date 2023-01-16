package test

import (
	"fmt"
	"path/filepath"

	"github.com/git-town/git-town/v7/test/helpers"
)

// GitManager manages the Git setup for the entire test suite.
// For each scenario, it provides a standardized, empty GitEnvironment consisting of a local and remote Git repository.
//
// Setting up a GitEnvironment is an expensive operation and has to be done for every scenario.
// As a performance optimization, GitManager creates a fully set up GitEnvironment (including the main branch and configuration)
// (the "memoized" environment) at the beginning of the test suite and makes copies of it for each scenario.
// Making copies of a fully set up Git repo is much faster than creating it from scratch.
// End-to-end tests run multi-threaded, all threads share a global GitManager instance.
type GitManager struct {
	counter helpers.Counter

	// path of the folder that this class operates in
	dir string

	// the memoized environment
	memoized GitEnvironment
}

// NewGitManager provides a new GitManager instance operating in the given directory.
func NewGitManager(dir string) (GitManager, error) {
	memoized, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
	if err != nil {
		return GitManager{}, fmt.Errorf("cannot create memoized environment: %w", err)
	}
	return GitManager{
		counter:  helpers.Counter{},
		dir:      dir,
		memoized: memoized,
	}, nil
}

// CreateScenarioEnvironment provides a new GitEnvironment for the scenario with the given name.
func (manager *GitManager) CreateScenarioEnvironment(scenarioName string) (GitEnvironment, error) {
	envDirName := helpers.FolderName(scenarioName) + "_" + manager.counter.ToString()
	envPath := filepath.Join(manager.dir, envDirName)
	return CloneGitEnvironment(manager.memoized, envPath)
}
