package infra

import (
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
	err := CopyDirectory(original.dir, dir)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot clone GitEnvironment '%s' to folder '%s'", original.dir, dir)
	}
	result := &GitEnvironment{
		dir:           dir,
		DeveloperRepo: LoadGitRepository(path.Join(dir, "developer")),
		OriginRepo:    LoadGitRepository(path.Join(dir, "origin")),
	}
	return result, nil
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
