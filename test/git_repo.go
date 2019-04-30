package test

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// GitRepository is a Git repository that exists inside a Git environment.
type GitRepository struct {
	dir string
	Runner
}

// InitGitRepository initializes a new Git repository in the given folder.
func InitGitRepository(dir string, bare bool) (GitRepository, error) {

	// create the folder
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot create directory %s", dir)
	}

	// cd into the folder
	err = os.Chdir(dir)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot cd into dir %s", dir)
	}

	// initialize the repo in the folder
	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}
	runner := Runner{}
	result := runner.Run("git", args...)
	if result.Err != nil {
		return GitRepository{}, errors.Wrap(result.Err, "error running git "+strings.Join(args, " "))
	}
	return GitRepository{dir: dir}, nil
}

// CloneFrom initializes this repository as a clone of the given parent repo.
func CloneFrom(parentDir, childDir string) (GitRepository, error) {

	// clone the repo
	runner := Runner{}
	result := runner.Run("git", "clone", parentDir, childDir)
	if result.Err != nil {
		return GitRepository{}, errors.Wrapf(result.Err, "cannot clone repo %s", parentDir)
	}

	// configure the repo
	err := os.Chdir(childDir)
	if err != nil {
		return GitRepository{}, err
	}
	userName := strings.Replace(path.Base(childDir), "_secondary", "", 1)
	err = runner.RunMany([][]string{
		[]string{"git", "config", "user.name", userName},
		[]string{"git", "config", "user.email", userName + "@example.com"},
		[]string{"git", "config", "push.default", "simple"},
		[]string{"git", "config", "core.editor", "vim"},
		[]string{"git", "config", "git-town.main-branch-name", "main"},
		[]string{"git", "config", "git-town.perennial-branch-names", ""},
	})
	return GitRepository{dir: childDir}, err
}
