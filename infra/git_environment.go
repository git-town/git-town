package infra

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
)

// GitEnvironment manages the setup, operation, and cleanup of a Git environment for a scenario.
type GitEnvironment struct {

	// dir is the directory that this environment is in
	dir string

	OriginRepo    *GitRepository
	DeveloperRepo *GitRepository
}

// NewGitEnvironment creates a new Git environment in the given directory
func NewGitEnvironment(baseDir string) (*GitEnvironment, error) {
	err := os.MkdirAll(baseDir, 0777)
	return &GitEnvironment{dir: baseDir}, err
}

// CloneFromFolder creates a new Git environment in the given folder
// that contains a copy of the given original environment.
func CloneEnvironment(original *GitEnvironment, dir string) (*GitEnvironment, error) {

	// create the folder for the new environment
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot make folder for cloned scenario: %s", dir)
	}

	// copy the memoized environment over
	runner := Runner{}
	tarCmd := fmt.Sprintf("(cd %s; tar c *) | tar xp", original.dir)
	fmt.Println("TAR", tarCmd)
	runResult := runner.Run("/bin/bash", "-c", tarCmd)
	if runResult.Err != nil {
		return nil, errors.Wrapf(runResult.Err, "cannot copy memoized environment over: %s", runResult.Output)
	}

	gitEnv := &GitEnvironment{dir: dir}
	gitEnv.DeveloperRepo = LoadExistingFolder(original.DeveloperRepo.dir)
	gitEnv.OriginRepo = LoadExistingFolder(path.Join(dir, "origin"))
	return gitEnv, nil
}

// CreateScenarioSetup creates the setup that all Cucumber Scenarios start out with.
func (ge GitEnvironment) CreateScenarioSetup() error {

	// Create origin repo and set "main" as default branch
	var err error
	ge.OriginRepo, err = InitGitRepository(ge.repositoryPath("origin"), true)
	if err != nil {
		return errors.Wrap(err, "cannot initialize origin directory")
	}

	// set "main" as default branch
	ge.OriginRepo.Run("git", "symbolic-ref", "HEAD", "refs/heads/main")

	// clone the "developer" repo
	ge.DeveloperRepo, err = CloneFrom(ge.repositoryPath("origin"), ge.repositoryPath("developer"))
	if err != nil {
		return errors.Wrap(err, "cannot clone developer repo from origin")
	}

	// Initialize main branch
	err = ge.DeveloperRepo.RunMany([][]string{
		[]string{"git", "checkout", "--orphan", "main"},
		[]string{"git", "commit", "--allow-empty", "-m", "Initial commit"},
		[]string{"git", "push", "-u", "origin", "main"},
	})
	if err != nil {
		return err
	}

	return nil
}

func (ge GitEnvironment) repositoryPath(repoName string) string {
	return path.Join(ge.dir, repoName)
}
