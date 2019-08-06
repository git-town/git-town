package infra

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
)

// GitEnvironment manages the Git environment for one test scenario.
type GitEnvironment struct {

	// dir is the directory that this environment is in
	dir string

	OriginRepo    GitRepository
	DeveloperRepo GitRepository
}

// NewGitEnvironment creates a new Git environment in the given directory
func NewGitEnvironment(baseDir string) (*GitEnvironment, error) {
	err := os.MkdirAll(baseDir, 0777)
	return &GitEnvironment{dir: baseDir}, err
}

// CloneGitEnvironment creates a new GitEnvironment in the given folder containing a copy of the given GitEnvironment.
func CloneGitEnvironment(original *GitEnvironment, dir string) (*GitEnvironment, error) {

	// create the folder for the new environment
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot make folder for scenario '%s'", dir)
	}

	// copy the folder contents of the memoized environment over
	runner := ShellRunner{}
	tarCmd := fmt.Sprintf("(cd %s; tar c *) | tar xp", original.dir)
	fmt.Println("TAR: ", tarCmd)
	runResult := runner.Run("/bin/bash", "-c", tarCmd)
	fmt.Println("TAR RESULT:", runResult)
	if runResult.Err != nil {
		return nil, errors.Wrapf(runResult.Err, "cannot copy memoized environment over: %s", runResult.Output)
	}
	gitEnv := &GitEnvironment{dir: dir}
	gitEnv.DeveloperRepo = LoadGitRepository(original.DeveloperRepo.dir)
	gitEnv.OriginRepo = LoadGitRepository(path.Join(dir, "origin"))
	return gitEnv, nil
}

// Populate instantiates the underlying folder content so that this GitEnvironment is ready for action.
// The name "populate" indicates that this takes a while.
func (ge GitEnvironment) Populate() error {

	// Create the origin repo and set "main" as the default branch
	var err error
	ge.OriginRepo, err = InitGitRepository(ge.repositoryPath("origin"), true)
	if err != nil {
		return errors.Wrap(err, "cannot initialize origin directory")
	}

	// set "main" as default branch
	ge.OriginRepo.Run("git", "symbolic-ref", "HEAD", "refs/heads/main")

	// git-clone the "developer" repo
	ge.DeveloperRepo, err = CloneGitRepository(ge.repositoryPath("origin"), ge.repositoryPath("developer"))
	if err != nil {
		return errors.Wrap(err, "cannot clone developer repo from origin")
	}

	// Initialize the main branch
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

// repositoryPath returns the full path to the Git repository with the given name.
func (ge GitEnvironment) repositoryPath(repoName string) string {
	return path.Join(ge.dir, repoName)
}
