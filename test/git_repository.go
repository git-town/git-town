package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
	"github.com/pkg/errors"
)

// GitRepository is a Git repository that exists inside a Git environment.
type GitRepository struct {

	// Dir contains the path of the directory that this repository is in.
	Dir string

	// originalCommits contains the commits in this repository before the test ran.
	originalCommits []Commit

	// ShellRunner enables to run console commands in this repo.
	ShellRunner

	// configCache contains the Git Town configuration to use.
	// This value is lazy loaded. Please use Configuration() to access it.
	configCache *git.Configuration
}

// NewGitRepository provides a new GitRepository instance working in the given directory.
// The directory must contain an existing Git repo.
func NewGitRepository(dir string, globalDir string) GitRepository {
	result := GitRepository{Dir: dir}
	result.ShellRunner = NewShellRunner(dir, globalDir)
	return result
}

// InitGitRepository initializes a fully functioning Git repository in the given path,
// including necessary Git configuration.
// Creates missing folders as needed.
func InitGitRepository(dir string, globalDir string) (GitRepository, error) {
	// create the folder
	err := os.MkdirAll(dir, 0744)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot create directory %q", dir)
	}

	// initialize the repo in the folder
	result := NewGitRepository(dir, globalDir)
	output, err := result.Run("git", "init")
	if err != nil {
		return result, errors.Wrapf(err, `error running "git init" in %q: %s`, dir, output)
	}
	err = result.RunMany([][]string{
		{"git", "config", "--global", "user.name", "user"},
		{"git", "config", "--global", "user.email", "email@example.com"},
		{"git", "config", "--global", "core.editor", "vim"},
	})
	return result, err
}

// CloneGitRepository clones the given parent repo into a new GitRepository.
func CloneGitRepository(parentDir, childDir, globalDir string) (GitRepository, error) {
	runner := NewShellRunner(".", globalDir)
	_, err := runner.Run("git", "clone", parentDir, childDir)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot clone repo %q", parentDir)
	}
	return NewGitRepository(childDir, globalDir), nil
}

// Branches provides the names of the local branches in this Git repository,
// sorted alphabetically.
func (repo *GitRepository) Branches() (result []string, err error) {
	output, err := repo.Run("git", "branch")
	if err != nil {
		return result, errors.Wrapf(err, "cannot run 'git branch' in repo %q", repo.dir)
	}
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.Replace(line, "* ", "", 1)
		result = append(result, strings.TrimSpace(line))
	}
	return sort.StringSlice(result), nil
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (repo *GitRepository) CheckoutBranch(name string) error {
	output, err := repo.Run("git", "checkout", name)
	if err != nil {
		return errors.Wrapf(err, "cannot check out branch %q in repo %q: %s", name, repo.dir, output)
	}
	return nil
}

// Commits provides a tabular list of the commits in this Git repository with the given fields.
func (repo *GitRepository) Commits(fields []string) (result []Commit, err error) {
	branches, err := repo.Branches()
	if err != nil {
		return result, errors.Wrap(err, "cannot determine the Git branches")
	}
	for _, branch := range branches {
		commits, err := repo.commitsInBranch(branch, fields)
		if err != nil {
			return result, err
		}
		result = append(result, commits...)
	}
	return result, nil
}

// CommitsInBranch provides all commits in the given Git branch.
func (repo *GitRepository) commitsInBranch(branch string, fields []string) (result []Commit, err error) {
	output, err := repo.Run("git", "log", branch, "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	if err != nil {
		return result, errors.Wrapf(err, "cannot get commits in branch %q", branch)
	}
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		parts := strings.Split(line, "|")
		commit := Commit{Branch: branch, SHA: parts[0], Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if util.DoesStringArrayContain(fields, "FILE NAME") {
			filenames, err := repo.FilesInCommit(commit.SHA)
			if err != nil {
				return result, errors.Wrapf(err, "cannot determine file name for commit %q in branch %q", commit.SHA, branch)
			}
			commit.FileName = strings.Join(filenames, ", ")
		}
		if util.DoesStringArrayContain(fields, "FILE CONTENT") {
			filecontent, err := repo.FileContentInCommit(commit.SHA, commit.FileName)
			if err != nil {
				return result, errors.Wrapf(err, "cannot determine file content for commit %q in branch %q", commit.SHA, branch)
			}
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result, nil
}

// Configuration lazy-loads the Git-Town configuration for this repo.
func (repo *GitRepository) Configuration() *git.Configuration {
	if repo.configCache == nil {
		repo.configCache = git.NewConfiguration(repo.Dir)
	}
	return repo.configCache
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (repo *GitRepository) CreateBranch(name string) error {
	output, err := repo.Run("git", "checkout", "-b", name)
	if err != nil {
		return errors.Wrapf(err, "cannot create branch %q: %s", name, output)
	}
	return nil
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (repo *GitRepository) CreateCommit(commit Commit) error {
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
	output, err = repo.Run("git", "commit", "-m", commit.Message)
	if err != nil {
		return errors.Wrapf(err, "cannot commit: %s", output)
	}
	return nil
}

// CreateFeatureBranch creates a branch with the given name in this repository.
func (repo *GitRepository) CreateFeatureBranch(name string) error {
	output, err := repo.Run("git", "town", "hack", name)
	if err != nil {
		return errors.Wrapf(err, "cannot create branch %q in repo: %s", name, output)
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

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func (repo *GitRepository) CreatePerennialBranches(names ...string) error {
	for _, name := range names {
		err := repo.CreateBranch(name)
		if err != nil {
			return errors.Wrapf(err, "cannot create perennial branch %q in repo %q", name, repo.Dir)
		}
	}
	repo.Configuration().AddToPerennialBranches(names...)
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

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (repo *GitRepository) FileContentInCommit(sha string, filename string) (result string, err error) {
	output, err := repo.Run("git", "show", sha+":"+filename)
	if err != nil {
		return result, errors.Wrapf(err, "cannot determine the content for file %q in commit %q", filename, sha)
	}
	return output, nil
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (repo *GitRepository) FilesInCommit(sha string) (result []string, err error) {
	output, err := repo.Run("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return result, errors.Wrapf(err, "cannot get files for commit %q", sha)
	}
	return strings.Split(strings.TrimSpace(output), "\n"), nil
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

// PushBranch pushes the branch with the given name to the remote.
func (repo *GitRepository) PushBranch(name string) error {
	output, err := repo.Run("git", "push", "-u", "origin", name)
	if err != nil {
		return errors.Wrapf(err, "cannot push branch %q in repo %q to origin: %s", name, repo.Dir, output)
	}
	return nil
}

// RegisterOriginalCommit tracks the given commit as existing in this repo before the system under test executed.
func (repo *GitRepository) RegisterOriginalCommit(commit Commit) {
	repo.originalCommits = append(repo.originalCommits, commit)
}

// SetOffline enables or disables offline mode for this GitRepository.
func (repo *GitRepository) SetOffline(enabled bool) error {
	output, err := repo.Run("git", "config", "--global", "git-town.offline", "true")
	if err != nil {
		return errors.Wrapf(err, "cannot set offline mode in repo %q: %s", repo.Dir, output)
	}
	return nil
}

// SetRemote sets the remote of this Git repository to the given target.
func (repo *GitRepository) SetRemote(target string) error {
	return repo.RunMany([][]string{
		{"git", "remote", "remove", "origin"},
		{"git", "remote", "add", "origin", target},
	})
}
