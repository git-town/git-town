package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Originate/git-town/test/gherkintools"
	"github.com/pkg/errors"
)

// GitRepository is a Git repository that exists inside a Git environment.
type GitRepository struct {

	// Dir contains the path of the directory that this repository is in.
	Dir string

	// originalCommits contains the commits in this repository before the test ran.
	originalCommits []gherkintools.Commit

	// ShellRunner enables to run console commands in this repo.
	ShellRunner
}

// NewGitRepository provides a new GitRepository instance working in the given directory.
// The directory must contain an existing Git repo.
func NewGitRepository(dir string) GitRepository {
	result := GitRepository{Dir: dir}
	result.ShellRunner = NewShellRunner(dir)
	return result
}

// InitGitRepository initializes a new Git repository in the given path.
// Creates missing folders as needed.
func InitGitRepository(dir string) (GitRepository, error) {
	// create the folder
	err := os.MkdirAll(dir, 0744)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot create directory %q", dir)
	}

	// initialize the repo in the folder
	result := NewGitRepository(dir)
	output, err := result.Run("git", "init")
	if err != nil {
		return result, errors.Wrapf(err, "error running git init in %q: %s", dir, output)
	}
	return result, nil
}

// CloneGitRepository clones the given parent repo into a new GitRepository.
func CloneGitRepository(parentDir, childDir string) (GitRepository, error) {
	runner := NewShellRunner(".")
	_, err := runner.Run("git", "clone", parentDir, childDir)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot clone repo %q", parentDir)
	}
	result := NewGitRepository(childDir)
	userName := strings.Replace(path.Base(childDir), "_secondary", "", 1)
	err = result.RunMany([][]string{
		{"git", "config", "user.name", userName},
		{"git", "config", "user.email", userName + "@example.com"},
		{"git", "config", "push.default", "simple"},
		{"git", "config", "core.editor", "vim"},
		{"git", "config", "git-town.main-branch-name", "main"},
		{"git", "config", "git-town.perennial-branch-names", ""},
	})
	return result, err
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (repo *GitRepository) CheckoutBranch(name string) error {
	output, err := repo.Run("git", "checkout", name)
	if err != nil {
		return errors.Wrapf(err, "cannot check out branch %q in repo %q: %s", name, repo.dir, output)
	}
	return nil
}

// createCommit creates a commit with the given properties in this Git repo.
func (repo *GitRepository) createCommit(commit gherkintools.Commit) error {
	repo.originalCommits = append(repo.originalCommits, commit)
	err := repo.CheckoutBranch(commit.Branch)
	if err != nil {
		return errors.Wrapf(err, "cannot checkout branch %q", commit.Branch)
	}
	err = repo.CreateFile(commit.FileName, commit.FileContent)
	if err != nil {
		return errors.Wrapf(err, "cannot create file %q needed for commit", commit.FileName)
	}
	output, err := repo.Run("git", "add", commit.FileName)
	if err != nil {
		return errors.Wrapf(err, "cannot add file to commit: %s", output)
	}
	_, err = repo.Run("git", "commit", "-m", commit.Message)
	if err != nil {
		return errors.Wrapf(err, "cannot commit: %s", output)
	}
	return nil
}

// CreateFile creates a file with the given name and content in this repository.
func (repo *GitRepository) CreateFile(name, content string) error {
	err := ioutil.WriteFile(path.Join(repo.Dir, name), []byte(content), 0744)
	if err != nil {
		return errors.Wrapf(err, "cannot create file %q", name)
	}
	return nil
}

// CurrentBranch provides the currently checked out branch for this repo.
func (repo *GitRepository) CurrentBranch() (result string, err error) {
	output, err := repo.Run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return result, errors.Wrapf(err, "cannot determine the current branch: %s", output)
	}
	return strings.TrimSpace(output), nil
}

// HasFile indicates whether this repository contains a file with the given name and content.
func (repo *GitRepository) HasFile(name, content string) (result bool, err error) {
	rawContent, err := ioutil.ReadFile(path.Join(repo.Dir, name))
	if err != nil {
		return result, errors.Wrapf(err, "repo doesn't have file %q", name)
	}
	actualContent := string(rawContent)
	if actualContent != content {
		return result, fmt.Errorf("file %q should have content %q but has %q", name, content, actualContent)
	}
	return true, nil
}

// SetRemote sets the remote of this Git repository to the given target.
func (repo *GitRepository) SetRemote(target string) error {
	return repo.RunMany([][]string{
		{"git", "remote", "remove", "origin"},
		{"git", "remote", "add", "origin", target},
	})
}
