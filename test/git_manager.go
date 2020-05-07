package test

import (
	"fmt"
	"path/filepath"

	"github.com/git-town/git-town/test/helpers"
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

// NewGitManager provides a new GitManager instance operating in the given directory.
func NewGitManager(baseDir string) *GitManager {
	return &GitManager{dir: baseDir}
}

// CreateMemoizedEnvironment creates the Git environment cache
// that makes cloning new GitEnvironment instances faster.
func (manager *GitManager) CreateMemoizedEnvironment() error {
	var err error
	manager.memoized, err = NewStandardGitEnvironment(filepath.Join(manager.dir, "memoized"))
	if err != nil {
		return fmt.Errorf("cannot create memoized environment: %w", err)
	}
	return nil
}

// CreateScenarioEnvironment provides a new GitEnvironment for the scenario with the given name
func (manager *GitManager) CreateScenarioEnvironment(scenarioName string) (*GitEnvironment, error) {
	envDirName := helpers.FolderName(scenarioName) + "_" + helpers.UniqueString()
	envPath := filepath.Join(manager.dir, envDirName)
	return CloneGitEnvironment(manager.memoized, envPath)
}
