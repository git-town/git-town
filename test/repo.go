package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/stretchr/testify/assert"
)

// TestRepo provides Git functionality for test code (unit and end-to-end tests).
type TestRepo struct {
	TestCommands
	internal git.InternalCommands
}

// CreateRepo creates TestRepo instances.
func CreateRepo(t *testing.T) TestRepo {
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
func InitRepo(workingDir, homeDir, binDir string) (TestRepo, error) {
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
// TODO: inline this method.
func NewRepo(workingDir, homeDir, binDir string) TestRepo {
	mockingRunner := MockingRunner{
		workingDir: workingDir,
		homeDir:    homeDir,
		binDir:     binDir,
	}
	config := git.RepoConfig{
		GitTown:            config.NewGitTown(&mockingRunner),
		CurrentBranchCache: &cache.String{},
		DryRun:             false,
		IsRepoCache:        &cache.Bool{},
		RemoteBranchCache:  &cache.Strings{},
		RemotesCache:       &cache.Strings{},
		RootDirCache:       &cache.String{},
	}
	internalCommands := git.InternalCommands{
		InternalRunner: git.NewInternalRunner(false),
		Config:         &config,
	}
	testCommands := TestCommands{
		MockingRunner:    mockingRunner,
		config:           config,
		InternalCommands: &internalCommands,
	}
	return TestRepo{
		TestCommands: testCommands,
		internal:     internalCommands,
	}
}

// CreateTestGitTownRepo creates a GitRepo for use in tests, with a main branch and
// initial git town configuration.
func CreateTestGitTownRepo(t *testing.T) TestRepo {
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
