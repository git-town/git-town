package test

import (
	"os"
	"path"

	"github.com/pkg/errors"
)

// GitEnvironment is the complete Git environment for a test scenario.
type GitEnvironment struct {

	// Dir is the directory that this environment is in.
	Dir string

	// OriginRepo is the Git repository that simulates the remote repo (on GitHub).
	OriginRepo GitRepository

	// DeveloperRepo is the Git repository that is locally checked out at the developer machine.
	DeveloperRepo GitRepository
}

// CloneGitEnvironment provides a GitEnvironment instance in the given directory,
// containing a copy of the given GitEnvironment.
func CloneGitEnvironment(original *GitEnvironment, dir string) (*GitEnvironment, error) {
	err := CopyDirectory(original.Dir, dir)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot clone GitEnvironment %q to folder %q", original.Dir, dir)
	}
	result := GitEnvironment{
		Dir:           dir,
		DeveloperRepo: NewGitRepository(path.Join(dir, "developer")),
		OriginRepo:    NewGitRepository(path.Join(dir, "origin")),
	}
	return &result, nil
}

// NewGitEnvironment provides a Git environment instance located in the given directory path.
// Missing directories are created as needed.
func NewGitEnvironment(baseDir string) (*GitEnvironment, error) {
	err := os.MkdirAll(baseDir, 0744)
	return &GitEnvironment{Dir: baseDir}, err
}

// NewStandardGitEnvironment provides a GitEnvironment in the given directory,
// fully populated as a standardized setup for scenarios.
func NewStandardGitEnvironment(dir string) (result *GitEnvironment, err error) {
	result, err = NewGitEnvironment(dir)
	if err != nil {
		return result, errors.Wrapf(err, "cannot create a new standard environment")
	}

	// create the origin repo
	result.OriginRepo, err = InitGitRepository(result.originRepoPath(), true)
	if err != nil {
		return result, errors.Wrapf(err, "cannot initialize origin directory at %q", result.originRepoPath())
	}

	// set "main" as the default branch
	_, err = result.OriginRepo.Run("git", "symbolic-ref", "HEAD", "refs/heads/main")
	if err != nil {
		return result, errors.Wrap(err, "cannot set main as the default branch")
	}

	// git-clone the "developer" repo
	result.DeveloperRepo, err = CloneGitRepository(result.originRepoPath(), result.developerRepoPath())
	if err != nil {
		return result, errors.Wrapf(err, "cannot clone developer repo (%q) from origin (%q)", result.originRepoPath(), result.developerRepoPath())
	}

	// initialize the main branch
	err = result.DeveloperRepo.RunMany([][]string{
		{"git", "checkout", "--orphan", "main"},
		{"git", "commit", "--allow-empty", "-m", "Initial commit"},
		{"git", "push", "-u", "origin", "main"},
	})
	return result, err
}

// Remove deletes all files used by this GitEnvironment from disk.
func (env GitEnvironment) Remove() error {
	return os.RemoveAll(env.Dir)
}

// developerRepoPath provides the full path to the Git repository with the given name.
func (env GitEnvironment) developerRepoPath() string {
	return path.Join(env.Dir, "developer")
}

// originRepoPath provides the full path to the Git repository with the given name.
func (env GitEnvironment) originRepoPath() string {
	return path.Join(env.Dir, "origin")
}
