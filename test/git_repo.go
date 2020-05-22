package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/util"
	"github.com/stretchr/testify/assert"
)

// GitRepo extends git.Repo with facilities for testing.
type GitRepo struct {
	*git.Repo
}

// CloneGitRepo clones a Git repo in originDir into a new GitRepository in workingDir.
// The cloning operation is using the given homeDir as the $HOME.
func CloneGitRepo(originDir, targetDir, homeDir, binDir string) (GitRepo, error) {
	res, err := command.Run("git", "clone", originDir, targetDir)
	if err != nil {
		return GitRepo{}, fmt.Errorf("cannot clone repo %q: %w\n%s", originDir, err, res.Output())
	}
	return NewGitRepository(targetDir, homeDir, NewMockingShell(targetDir, homeDir, binDir)), nil
}

// CreateTestRepo creates a GitRepo for use in tests
func CreateTestRepo(t *testing.T) GitRepo {
	dir := CreateTempDir(t)
	repo, err := InitGitRepository(dir, dir, "")
	assert.Nil(t, err, "cannot initialize Git repo")
	err = repo.Shell.RunMany([][]string{
		{"git", "commit", "--allow-empty", "-m", "initial commit"},
	})
	assert.Nil(t, err, "cannot create initial commit: %s")
	return repo
}

// InitGitRepository initializes a fully functioning Git repository in the given path,
// including necessary Git configuration.
// Creates missing folders as needed.
func InitGitRepository(workingDir, homeDir, binDir string) (GitRepo, error) {
	// create the folder
	err := os.MkdirAll(workingDir, 0744)
	if err != nil {
		return GitRepo{}, fmt.Errorf("cannot create directory %q: %w", workingDir, err)
	}
	// initialize the repo in the folder
	result := NewGitRepository(workingDir, homeDir, NewMockingShell(workingDir, homeDir, binDir))
	outcome, err := result.Shell.Run("git", "init")
	if err != nil {
		return result, fmt.Errorf(`error running "git init" in %q: %w\n%v`, workingDir, err, outcome)
	}
	err = result.Shell.RunMany([][]string{
		{"git", "config", "--global", "user.name", "user"},
		{"git", "config", "--global", "user.email", "email@example.com"},
		{"git", "config", "--global", "core.editor", "vim"},
	})
	return result, err
}

// NewGitRepository provides a new GitRepository instance working in the given directory.
// The directory must contain an existing Git repo.
// TODO: remove homeDir here, it is included in the given Shell.
func NewGitRepository(workingDir string, homeDir string, shell command.Shell) GitRepo {
	repo := git.Repo{Dir: workingDir, Shell: shell}
	return GitRepo{&repo}
}

// FilesInBranches provides a data table of files and their content in all branches.
func (repo *GitRepo) FilesInBranches() (result DataTable, err error) {
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := repo.Branches()
	if err != nil {
		return result, err
	}
	for _, branch := range branches {
		files, err := repo.FilesInBranch(branch)
		if err != nil {
			return result, err
		}
		for _, file := range files {
			content, err := repo.FileContentInCommit(branch, file)
			if err != nil {
				return result, err
			}
			result.AddRow(branch, file, content)
		}
	}
	return result, err
}

// Commits provides a tabular list of the commits in this Git repository with the given fields.
func (repo *GitRepo) Commits(fields []string) (result []Commit, err error) {
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
func (repo *GitRepo) commitsInBranch(branch string, fields []string) (result []Commit, err error) {
	outcome, err := repo.Shell.Run("git", "log", branch, "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	if err != nil {
		return result, fmt.Errorf("cannot get commits in branch %q: %w", branch, err)
	}
	for _, line := range strings.Split(outcome.OutputSanitized(), "\n") {
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

// CreateCommit creates a commit with the given properties in this Git repo.
func (repo *GitRepo) CreateCommit(commit Commit) error {
	err := repo.CheckoutBranch(commit.Branch)
	if err != nil {
		return fmt.Errorf("cannot checkout branch %q: %w", commit.Branch, err)
	}
	err = repo.CreateFile(commit.FileName, commit.FileContent)
	if err != nil {
		return fmt.Errorf("cannot create file %q needed for commit: %w", commit.FileName, err)
	}
	outcome, err := repo.Shell.Run("git", "add", commit.FileName)
	if err != nil {
		return fmt.Errorf("cannot add file to commit: %w\n%v", err, outcome)
	}
	commands := []string{"commit", "-m", commit.Message}
	if commit.Author != "" {
		commands = append(commands, "--author="+commit.Author)
	}
	outcome, err = repo.Shell.Run("git", commands...)
	if err != nil {
		return fmt.Errorf("cannot commit: %w\n%v", err, outcome)
	}
	return nil
}
