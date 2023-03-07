package test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/run"
	"github.com/stretchr/testify/assert"
)

// Repo is a Git Repo in test code.
type Repo struct {
	git.Runner              // the git.Runner instance to use
	shell      MockingShell // a reference to the MockingShell instance used here
}

// CreateRepo creates TestRepo instances.
func CreateRepo(t *testing.T) Repo {
	t.Helper()
	dir := t.TempDir()
	workingDir := filepath.Join(dir, "repo")
	err := os.Mkdir(workingDir, 0o744)
	assert.NoError(t, err)
	homeDir := filepath.Join(dir, "home")
	err = os.Mkdir(homeDir, 0o744)
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
		{"git", "init", "--initial-branch=initial"},
		{"git", "config", "--global", "user.name", "user"},
		{"git", "config", "--global", "user.email", "email@example.com"},
	})
	return result, err
}

// NewRepo provides a new Repo instance working in the given directory.
// The directory must contain an existing Git repo.
func NewRepo(workingDir, homeDir, binDir string) Repo {
	shell := NewMockingShell(workingDir, homeDir, binDir)
	runner := git.Runner{
		Shell:              &shell,
		Config:             config.NewGitTown(&shell),
		DryRun:             &git.DryRun{},
		IsRepoCache:        &cache.Bool{},
		RemoteBranchCache:  &cache.Strings{},
		RemotesCache:       &cache.Strings{},
		RootDirCache:       &cache.String{},
		CurrentBranchCache: &cache.String{},
	}
	return Repo{Runner: runner, shell: shell}
}

// BranchHierarchyTable provides the currently configured branch hierarchy information as a DataTable.
func (repo *Repo) BranchHierarchyTable() DataTable {
	result := DataTable{}
	repo.Config.Reload()
	parentBranchMap := repo.Config.ParentBranchMap()
	result.AddRow("BRANCH", "PARENT")
	childBranches := make([]string, 0, len(parentBranchMap))
	for child := range parentBranchMap {
		childBranches = append(childBranches, child)
	}
	sort.Strings(childBranches)
	for _, child := range childBranches {
		result.AddRow(child, parentBranchMap[child])
	}
	return result
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
func (repo *Repo) FilesInBranches() (DataTable, error) {
	result := DataTable{}
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := repo.LocalBranchesMainFirst()
	if err != nil {
		return DataTable{}, err
	}
	lastBranch := ""
	for _, branch := range branches {
		files, err := repo.FilesInBranch(branch)
		if err != nil {
			return DataTable{}, err
		}
		for _, file := range files {
			content, err := repo.FileContentInCommit(branch, file)
			if err != nil {
				return DataTable{}, err
			}
			if branch == lastBranch {
				result.AddRow("", file, content)
			} else {
				result.AddRow(branch, file, content)
			}
			lastBranch = branch
		}
	}
	return result, err
}

// CreateTestGitTownRepo creates a GitRepo for use in tests, with a main branch and
// initial git town configuration.
func CreateTestGitTownRepo(t *testing.T) Repo {
	t.Helper()
	repo := CreateRepo(t)
	err := repo.CreateBranch("main", "initial")
	assert.NoError(t, err)
	err = repo.Config.SetMainBranch("main")
	assert.NoError(t, err)
	err = repo.Config.SetPerennialBranches([]string{})
	assert.NoError(t, err)
	return repo
}
