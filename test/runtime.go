package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v8/src/cache"
	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/src/execute"
	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/subshell"
	"github.com/git-town/git-town/v8/test/runner"
	"github.com/stretchr/testify/assert"
)

// Runtime provides Git functionality for test code (unit and end-to-end tests).
type Runtime struct {
	testCommands
	Backend git.BackendCommands
}

// CreateRuntime creates test.Runner instances.
func CreateRuntime(t *testing.T) Runtime {
	t.Helper()
	dir := t.TempDir()
	workingDir := filepath.Join(dir, "repo")
	err := os.Mkdir(workingDir, 0o744)
	assert.NoError(t, err)
	homeDir := filepath.Join(dir, "home")
	err = os.Mkdir(homeDir, 0o744)
	assert.NoError(t, err)
	runner, err := initRuntime(workingDir, homeDir, homeDir)
	assert.NoError(t, err)
	_, err = runner.Run("git", "commit", "--allow-empty", "-m", "initial commit")
	assert.NoError(t, err)
	return runner
}

// initRuntime creates a fully functioning test.Runner in the given working directory,
// including necessary Git configuration to make commits. Creates missing folders as needed.
func initRuntime(workingDir, homeDir, binDir string) (Runtime, error) {
	runner := newRuntime(workingDir, homeDir, binDir)
	err := runner.RunMany([][]string{
		{"git", "init", "--initial-branch=initial"},
		{"git", "config", "--global", "user.name", "user"},
		{"git", "config", "--global", "user.email", "email@example.com"},
	})
	return runner, err
}

// newRuntime provides a new test.Runner instance working in the given directory.
// The directory must contain an existing Git repo.
func newRuntime(workingDir, homeDir, binDir string) Runtime {
	mockingRunner := runner.Mocking{
		WorkingDir: workingDir,
		HomeDir:    homeDir,
		BinDir:     binDir,
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
		BackendRunner: subshell.BackendRunner{Dir: &workingDir, Verbose: false, Stats: &execute.NoStatistics{}},
		Config:        &config,
	}
	testCommands := testCommands{
		Mocking:         mockingRunner,
		config:          config,
		BackendCommands: &backendCommands,
	}
	return Runtime{
		testCommands: testCommands,
		Backend:      backendCommands,
	}
}

// CreateTestGitTownRuntime creates a test.Runtime for use in tests,
// with a main branch and initial git town configuration.
func CreateTestGitTownRuntime(t *testing.T) Runtime {
	t.Helper()
	repo := CreateRuntime(t)
	err := repo.CreateBranch("main", "initial")
	assert.NoError(t, err)
	err = repo.Config.SetMainBranch("main")
	assert.NoError(t, err)
	err = repo.Config.SetPerennialBranches([]string{})
	assert.NoError(t, err)
	return repo
}
