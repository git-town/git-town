package test

import (
	"os"
	"path"

	"github.com/pkg/errors"
)

// GitEnvironment is the complete Git environment for a test scenario.
type GitEnvironment struct {

	// dir is the directory that this environment is in.
	dir string

	// OriginRepo is the Git repository that simulates the remote repo (on GitHub).
	OriginRepo GitRepository

	// DeveloperRepo is the Git repository that is locally checked out at the developer machine.
	DeveloperRepo GitRepository
}

// NewGitEnvironment provides a Git environment instance located in the given directory path.
// Missing directories are created as needed.
func NewGitEnvironment(baseDir string) (*GitEnvironment, error) {
	err := os.MkdirAll(baseDir, 0777)
	return &GitEnvironment{dir: baseDir}, err
}

// CloneGitEnvironment provides a GitEnvironment instance in the given directory,
// containing a copy of the given GitEnvironment.
func CloneGitEnvironment(original *GitEnvironment, dir string) (*GitEnvironment, error) {
	err := CopyDirectory(original.dir, dir)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot clone GitEnvironment %q to folder %q", original.dir, dir)
	}
	result := GitEnvironment{
		dir:           dir,
		DeveloperRepo: NewGitRepository(path.Join(dir, "developer")),
		OriginRepo:    NewGitRepository(path.Join(dir, "origin")),
	}
	return &result, nil
}

// Populate instantiates the underlying folder content so that this GitEnvironment is ready for action.
// The name "populate" indicates that this takes a while.
func (ge GitEnvironment) Populate() error {

	// create the origin repo
	var err error
	ge.OriginRepo, err = InitGitRepository(ge.originRepoPath(), true)
	if err != nil {
		return errors.Wrapf(err, "cannot initialize origin directory at %q", ge.originRepoPath())
	}

	// set "main" as the default branch
	ge.OriginRepo.Run("git", "symbolic-ref", "HEAD", "refs/heads/main")

	// git-clone the "developer" repo
	ge.DeveloperRepo, err = CloneGitRepository(ge.originRepoPath(), ge.developerRepoPath())
	if err != nil {
		return errors.Wrapf(err, "cannot clone developer repo (%q) from origin (%q)", ge.originRepoPath(), ge.developerRepoPath())
	}

	// Initialize the main branch
	err = ge.DeveloperRepo.RunMany([][]string{
		[]string{"git", "checkout", "--orphan", "main"},
		[]string{"git", "commit", "--allow-empty", "-m", "Initial commit"},
		[]string{"git", "push", "-u", "origin", "main"},
	})
	return err
}

// developerRepoPath provides the full path to the Git repository with the given name.
func (ge GitEnvironment) developerRepoPath() string {
	return path.Join(ge.dir, "developer")
}

// originRepoPath provides the full path to the Git repository with the given name.
func (ge GitEnvironment) originRepoPath() string {
	return path.Join(ge.dir, "origin")
}
