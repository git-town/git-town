package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/execute"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/stretchr/testify/assert"
)

// Runner provides Git functionality for test code (unit and end-to-end tests).
type Runner struct {
	testCommands
	Backend git.BackendCommands
}

// CreateRunner creates test.Runner instances.
func CreateRunner(t *testing.T) Runner {
	t.Helper()
	dir := t.TempDir()
	workingDir := filepath.Join(dir, "repo")
	err := os.Mkdir(workingDir, 0o744)
	assert.NoError(t, err)
	homeDir := filepath.Join(dir, "home")
	err = os.Mkdir(homeDir, 0o744)
	assert.NoError(t, err)
	runner, err := initRunner(workingDir, homeDir, homeDir)
	assert.NoError(t, err)
	_, err = runner.Run("git", "commit", "--allow-empty", "-m", "initial commit")
	assert.NoError(t, err)
	return runner
}

// initRunner creates a fully functioning test.Runner in the given working directory,
// including necessary Git configuration to make commits. Creates missing folders as needed.
func initRunner(workingDir, homeDir, binDir string) (Runner, error) {
	result := newRunner(workingDir, homeDir, binDir)
	err := result.RunMany([][]string{
		{"git", "init", "--initial-branch=initial"},
		{"git", "config", "--global", "user.name", "user"},
		{"git", "config", "--global", "user.email", "email@example.com"},
	})
	return result, err
}

// newRunner provides a new test.Runner instance working in the given directory.
// The directory must contain an existing Git repo.
// TODO: inline this method.
func newRunner(workingDir, homeDir, binDir string) Runner {
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
	backendCommands := git.BackendCommands{
		BackendRunner: execute.NewBackendRunner(&workingDir, false, nil),
		Config:        &config,
	}
	testCommands := testCommands{
		MockingRunner:   mockingRunner,
		config:          config,
		BackendCommands: &backendCommands,
	}
	return Runner{
		testCommands: testCommands,
		Backend:      backendCommands,
	}
}

// CreateTestGitTownRunner creates a test.Runner for use in tests,
// with a main branch and initial git town configuration.
func CreateTestGitTownRunner(t *testing.T) Runner {
	t.Helper()
	repo := CreateRunner(t)
	err := repo.CreateBranch("main", "initial")
	assert.NoError(t, err)
	err = repo.Config.SetMainBranch("main")
	assert.NoError(t, err)
	err = repo.Config.SetPerennialBranches([]string{})
	assert.NoError(t, err)
	return repo
}
