package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/src/config"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/run"
	"github.com/stretchr/testify/assert"
)

// Repo is a Git Repo in test code.
type Repo struct {
	git.Runner               // the git.Runner instance to use
	shell      *MockingShell // a reference to the MockingShell instance used here
}

// CreateRepo creates TestRepo instances.
func CreateRepo(t *testing.T) Repo {
	dir := CreateTempDir(t)
	workingDir := filepath.Join(dir, "repo")
	err := os.Mkdir(workingDir, 0744)
	assert.NoError(t, err)
	homeDir := filepath.Join(dir, "home")
	err = os.Mkdir(homeDir, 0744)
	assert.NoError(t, err)
	repo, err := InitRepo(workingDir, homeDir, homeDir)
	assert.NoError(t, err)
	_, err = repo.Run("git", "commit", "--allow-empty", "-m", "initial commit")
	assert.NoError(t, err)
	return repo
}

// InitRepo creates a fully functioning test.Repo in the given working directory,
// including necessary Git configuration to make commits. Creates missing folders as needed.
func InitRepo(workingDir, homeDir, binDir string) (Repo, error) {
	result := NewRepo(workingDir, homeDir, binDir)
	err := result.RunMany([][]string{
		{"git", "init"},
		{"git", "config", "--global", "user.name", "user"},
		{"git", "config", "--global", "user.email", "email@example.com"},
		{"git", "config", "--global", "core.editor", "vim"},
	})
	return result, err
}

// NewRepo provides a new Repo instance working in the given directory.
// The directory must contain an existing Git repo.
func NewRepo(workingDir, homeDir, binDir string) Repo {
	shell := NewMockingShell(workingDir, homeDir, binDir)
	runner := git.Runner{
		Shell:              shell,
		Config:             config.NewConfiguration(shell),
		DryRun:             &git.DryRun{},
		IsRepoCache:        &git.BoolCache{},
		RemoteBranchCache:  &git.StringSliceCache{},
		RemotesCache:       &git.StringSliceCache{},
		RootDirCache:       &git.StringCache{},
		CurrentBranchCache: &git.StringCache{},
	}
	return Repo{runner, shell}
}

// Clone creates a clone of this Repo into the given directory.
// The cloned repo uses the same homeDir and binDir as its origin.
func (repo *Repo) Clone(targetDir string) (Repo, error) {
	_, err := run.Exec("git", "clone", repo.shell.workingDir, targetDir)
	if err != nil {
		return Repo{}, fmt.Errorf("cannot clone repo %q: %w", repo.shell.workingDir, err)
	}
	return NewRepo(targetDir, repo.shell.homeDir, repo.shell.binDir), nil
}

// FilesInBranches provides a data table of files and their content in all branches.
func (repo *Repo) FilesInBranches() (result DataTable, err error) {
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := repo.LocalBranchesMainFirst()
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

// CreateTestGitTownRepo creates a GitRepo for use in tests, with a main branch and
// initial git town configuration.
func CreateTestGitTownRepo(t *testing.T) Repo {
	repo := CreateRepo(t)
	err := repo.CreateBranch("main", "master")
	assert.NoError(t, err)
	err = repo.RunMany([][]string{
		{"git", "config", "git-town.main-branch-name", "main"},
		{"git", "config", "git-town.perennial-branch-names", ""},
	})
	assert.NoError(t, err)
	return repo
}
