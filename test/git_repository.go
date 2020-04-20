package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
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

// CloneGitRepository clones a Git repo in originDir into a new GitRepository in workingDir.
// The cloning operation is using the given homeDir as the $HOME.
func CloneGitRepository(originDir, targetDir, homeDir string) (GitRepository, error) {
	runner := NewShellRunner(".", homeDir)
	res, err := runner.Run("git", "clone", originDir, targetDir)
	if err != nil {
		return GitRepository{}, fmt.Errorf("cannot clone repo %q: %w\n%s", originDir, err, res.Output())
	}
	return NewGitRepository(targetDir, homeDir), nil
}

// InitGitRepository initializes a fully functioning Git repository in the given path,
// including necessary Git configuration.
// Creates missing folders as needed.
func InitGitRepository(workingDir string, homeDir string) (GitRepository, error) {
	// create the folder
	err := os.MkdirAll(workingDir, 0744)
	if err != nil {
		return GitRepository{}, fmt.Errorf("cannot create directory %q: %w", workingDir, err)
	}
	// initialize the repo in the folder
	result := NewGitRepository(workingDir, homeDir)
	outcome, err := result.Run("git", "init")
	if err != nil {
		return result, fmt.Errorf(`error running "git init" in %q: %w\n%v`, workingDir, err, outcome)
	}
	err = result.RunMany([][]string{
		{"git", "config", "--global", "user.name", "user"},
		{"git", "config", "--global", "user.email", "email@example.com"},
		{"git", "config", "--global", "core.editor", "vim"},
	})
	return result, err
}

// NewGitRepository provides a new GitRepository instance working in the given directory.
// The directory must contain an existing Git repo.
func NewGitRepository(workingDir string, homeDir string) GitRepository {
	result := GitRepository{Dir: workingDir}
	result.ShellRunner = NewShellRunner(workingDir, homeDir)
	return result
}

// Branches provides the names of the local branches in this Git repository,
// sorted alphabetically.
func (repo *GitRepository) Branches() (result []string, err error) {
	outcome, err := repo.Run("git", "branch")
	if err != nil {
		return result, fmt.Errorf("cannot run 'git branch' in repo %q: %w", repo.workingDir, err)
	}
	for _, line := range strings.Split(strings.TrimSpace(outcome.OutputSanitized()), "\n") {
		line = strings.Replace(line, "* ", "", 1)
		result = append(result, strings.TrimSpace(line))
	}
	return sort.StringSlice(result), nil
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (repo *GitRepository) CheckoutBranch(name string) error {
	outcome, err := repo.Run("git", "checkout", name)
	if err != nil {
		return fmt.Errorf("cannot check out branch %q in repo %q: %w\n%v", name, repo.workingDir, err, outcome)
	}
	return nil
}

// Commits provides a tabular list of the commits in this Git repository with the given fields.
func (repo *GitRepository) Commits(fields []string) (result []Commit, err error) {
	branches, err := repo.Branches()
	if err != nil {
		return result, fmt.Errorf("cannot determine the Git branches: %w", err)
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
	outcome, err := repo.Run("git", "log", branch, "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	if err != nil {
		return result, fmt.Errorf("cannot get commits in branch %q: %w", branch, err)
	}
	for _, line := range strings.Split(strings.TrimSpace(outcome.OutputSanitized()), "\n") {
		parts := strings.Split(line, "|")
		commit := Commit{Branch: branch, SHA: parts[0], Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if util.DoesStringArrayContain(fields, "FILE NAME") {
			filenames, err := repo.FilesInCommit(commit.SHA)
			if err != nil {
				return result, fmt.Errorf("cannot determine file name for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileName = strings.Join(filenames, ", ")
		}
		if util.DoesStringArrayContain(fields, "FILE CONTENT") {
			filecontent, err := repo.FileContentInCommit(commit.SHA, commit.FileName)
			if err != nil {
				return result, fmt.Errorf("cannot determine file content for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result, nil
}

// CommitStagedChanges commits the currently staged changes.
func (repo *GitRepository) CommitStagedChanges(message bool) error {
	var out *command.Result
	var err error
	if message {
		out, err = repo.Run("git", "commit", "-m", "committing")
	} else {
		out, err = repo.Run("git", "commit", "--no-edit")
	}
	if err != nil {
		return fmt.Errorf("cannot commit staged changes: %w\n%s", err, out)
	}
	return nil
}

// Configuration returns a cached Configuration instance for this repo.
func (repo *GitRepository) Configuration(refresh bool) *git.Configuration {
	if repo.configCache == nil || refresh {
		repo.configCache = git.NewConfiguration(repo.Dir)
	}
	return repo.configCache
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (repo *GitRepository) CreateBranch(name string) error {
	outcome, err := repo.Run("git", "checkout", "-b", name)
	if err != nil {
		return fmt.Errorf("cannot create branch %q: %w\n%v", name, err, outcome)
	}
	return nil
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (repo *GitRepository) CreateChildFeatureBranch(name string, parentBranch string) error {
	err := repo.CheckoutBranch(parentBranch)
	if err != nil {
		return fmt.Errorf("cannot checkout parent branch %q: %w", parentBranch, err)
	}
	outcome, err := repo.Run("git", "town", "append", name)
	if err != nil {
		return fmt.Errorf("cannot create child branch %q: %w\n%v", name, err, outcome)
	}
	return nil
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (repo *GitRepository) CreateCommit(commit Commit) error {
	repo.originalCommits = append(repo.originalCommits, commit)
	err := repo.CheckoutBranch(commit.Branch)
	if err != nil {
		return fmt.Errorf("cannot checkout branch %q: %w", commit.Branch, err)
	}
	err = repo.CreateFile(commit.FileName, commit.FileContent)
	if err != nil {
		return fmt.Errorf("cannot create file %q needed for commit: %w", commit.FileName, err)
	}
	outcome, err := repo.Run("git", "add", commit.FileName)
	if err != nil {
		return fmt.Errorf("cannot add file to commit: %w\n%v", err, outcome)
	}
	outcome, err = repo.Run("git", "commit", "-m", commit.Message)
	if err != nil {
		return fmt.Errorf("cannot commit: %w\n%v", err, outcome)
	}
	return nil
}

// CreateFeatureBranch creates a branch with the given name in this repository.
func (repo *GitRepository) CreateFeatureBranch(name string) error {
	outcome, err := repo.Run("git", "town", "hack", name)
	if err != nil {
		return fmt.Errorf("cannot create branch %q in repo: %w\n%v", name, err, outcome)
	}
	return nil
}

// CreateFile creates a file with the given name and content in this repository.
func (repo *GitRepository) CreateFile(name, content string) error {
	filePath := filepath.Join(repo.Dir, name)
	folderPath := filepath.Dir(filePath)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create folder %q: %v", folderPath, err)
	}
	err = ioutil.WriteFile(filePath, []byte(content), 0744)
	if err != nil {
		return fmt.Errorf("cannot create file %q: %w", name, err)
	}
	return nil
}

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func (repo *GitRepository) CreatePerennialBranches(names ...string) error {
	for _, name := range names {
		err := repo.CreateBranch(name)
		if err != nil {
			return fmt.Errorf("cannot create perennial branch %q in repo %q: %w", name, repo.Dir, err)
		}
	}
	repo.Configuration(false).AddToPerennialBranches(names...)
	return nil
}

// CurrentBranch provides the currently checked out branch for this repo.
func (repo *GitRepository) CurrentBranch() (result string, err error) {
	outcome, err := repo.Run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return result, fmt.Errorf("cannot determine the current branch: %w\n%s", err, outcome.Output())
	}
	return strings.TrimSpace(outcome.OutputSanitized()), nil
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (repo *GitRepository) FileContentInCommit(sha string, filename string) (result string, err error) {
	outcome, err := repo.Run("git", "show", sha+":"+filename)
	if err != nil {
		return result, fmt.Errorf("cannot determine the content for file %q in commit %q: %w", filename, sha, err)
	}
	return outcome.OutputSanitized(), nil
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (repo *GitRepository) FilesInCommit(sha string) (result []string, err error) {
	outcome, err := repo.Run("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return result, fmt.Errorf("cannot get files for commit %q: %w", sha, err)
	}
	return strings.Split(strings.TrimSpace(outcome.OutputSanitized()), "\n"), nil
}

// HasFile indicates whether this repository contains a file with the given name and content.
func (repo *GitRepository) HasFile(name, content string) (result bool, err error) {
	rawContent, err := ioutil.ReadFile(filepath.Join(repo.Dir, name))
	if err != nil {
		return result, fmt.Errorf("repo doesn't have file %q: %w", name, err)
	}
	actualContent := string(rawContent)
	if actualContent != content {
		return result, fmt.Errorf("file %q should have content %q but has %q", name, content, actualContent)
	}
	return true, nil
}

// LastActiveDir provides the directory that was last used in this repo.
func (repo *GitRepository) LastActiveDir() (string, error) {
	res, err := repo.Run("git", "rev-parse", "--show-toplevel")
	return res.OutputSanitized(), err
}

// PushBranch pushes the branch with the given name to the remote.
func (repo *GitRepository) PushBranch(name string) error {
	outcome, err := repo.Run("git", "push", "-u", "origin", name)
	if err != nil {
		return fmt.Errorf("cannot push branch %q in repo %q to origin: %w\n%v", name, repo.Dir, err, outcome)
	}
	return nil
}

// RegisterOriginalCommit tracks the given commit as existing in this repo before the system under test executed.
func (repo *GitRepository) RegisterOriginalCommit(commit Commit) {
	repo.originalCommits = append(repo.originalCommits, commit)
}

// SetOffline enables or disables offline mode for this GitRepository.
func (repo *GitRepository) SetOffline(enabled bool) error {
	outcome, err := repo.Run("git", "config", "--global", "git-town.offline", "true")
	if err != nil {
		return fmt.Errorf("cannot set offline mode in repo %q: %w\n%v", repo.Dir, err, outcome)
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

// StageFiles adds the file with the given name to the Git index.
func (repo *GitRepository) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	_, err := repo.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot stage files %s: %w", strings.Join(names, ", "), err)
	}
	return nil
}
