package test

import (
	"path"

	"github.com/pkg/errors"
)

/*
GitManager manages the various Git repos needed to test Git operations.

The pre-defined environment in which scenarios start out in looks like this:
- the "developer" folder contains our local workspace repo
  (where we run tests in)
- the "origin" folder contains the remote repo for the developer repo
  (where pushes from "developer" go to)
- all repos contain a "main" branch that is configured as Git Town's main branch

Setting up the standardized environment happens a lot (before each scenario).
To make this process performant,
a fresh setup is cached in the "memoized" directory.
Before each scenario starts, it is copied into a folder for that scenario.

The folder structure is:
baseDir
├── memoized            # cache of the pre-defined environment for scenarios
|   ├── developer
|   └── origin
├── scenario A          # workspace for currently tested scenario A
|   ├── developer
|   └── origin
└── scenario B          # workspace for currently tested scenario B
    ├── developer
    └── origin
*/
type GitManager struct {

	// dir contains the name of the folder that this class operates in.
	dir string

	// the memoized environment
	memoized *GitEnvironment
}

// NewEnvironments creates a new Environments instance
// and prepopulates its environment cache.
func NewGitManager(baseDir string) *GitManager {
	return &GitManager{dir: baseDir}
}

// CreateMemoizedEnvironment creates the memoized environment
func (gm *GitManager) CreateMemoizedEnvironment() error {
	var err error
	gm.memoized, err = NewGitEnvironment(path.Join(gm.dir, "memoized"))
	if err != nil {
		return errors.Wrapf(err, "cannot create memoized environment")
	}
	err = gm.memoized.CreateScenarioSetup()
	if err != nil {
		return errors.Wrapf(err, "cannot populate memoized environment")
	}
	return nil
}

// CreateScenarioEnvironment creates a new GitEnvironment for the scenario with the given name
func (gm GitManager) CreateScenarioEnvironment(scenarioName string) (*GitEnvironment, error) {
	envPath := path.Join(gm.dir, scenarioName)
	return CloneEnvironment(envPath, gm.memoized)
}
