package testruntime

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v8/src/cache"
	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/src/execute"
	"github.com/git-town/git-town/v8/src/git"
	prodshell "github.com/git-town/git-town/v8/src/subshell"
	"github.com/git-town/git-town/v8/test/commands"
	testshell "github.com/git-town/git-town/v8/test/subshell"
	"github.com/stretchr/testify/assert"
)

// TestRuntime provides Git functionality for test code (unit and end-to-end tests).
type TestRuntime struct {
	commands.TestCommands
	Backend git.BackendCommands
}

// Create creates test.Runner instances.
func Create(t *testing.T) TestRuntime {
	t.Helper()
	dir := t.TempDir()
	workingDir := filepath.Join(dir, "repo")
	err := os.Mkdir(workingDir, 0o744)
	assert.NoError(t, err)
	homeDir := filepath.Join(dir, "home")
	err = os.Mkdir(homeDir, 0o744)
	assert.NoError(t, err)
	runtime, err := Initialize(workingDir, homeDir, homeDir)
	assert.NoError(t, err)
	_, err = runtime.Run("git", "commit", "--allow-empty", "-m", "initial commit")
	assert.NoError(t, err)
	return runtime
}

// initialize creates a fully functioning test.Runner in the given working directory,
// including necessary Git configuration to make commits. Creates missing folders as needed.
func Initialize(workingDir, homeDir, binDir string) (TestRuntime, error) {
	runtime := New(workingDir, homeDir, binDir)
	err := runtime.RunMany([][]string{
		{"git", "init", "--initial-branch=initial"},
		{"git", "config", "--global", "user.name", "user"},
		{"git", "config", "--global", "user.email", "email@example.com"},
	})
	return runtime, err
}

// newRuntime provides a new test.Runner instance working in the given directory.
// The directory must contain an existing Git repo.
func New(workingDir, homeDir, binDir string) TestRuntime {
	mockingRunner := testshell.Mocking{
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
		BackendRunner: prodshell.BackendRunner{Dir: &workingDir, Verbose: false, Stats: &execute.NoStatistics{}},
		Config:        &config,
	}
	testCommands := commands.TestCommands{
		Mocking:         mockingRunner,
		BackendCommands: &backendCommands,
	}
	return TestRuntime{
		TestCommands: testCommands,
		Backend:      backendCommands,
	}
}

// CreateGitTown creates a test.Runtime for use in tests,
// with a main branch and initial git town configuration.
func CreateGitTown(t *testing.T) TestRuntime {
	t.Helper()
	repo := Create(t)
	err := repo.CreateBranch("main", "initial")
	assert.NoError(t, err)
	err = repo.Config.SetMainBranch("main")
	assert.NoError(t, err)
	err = repo.Config.SetPerennialBranches([]string{})
	assert.NoError(t, err)
	return repo
}

// Clone creates a clone of the repository managed by this test.Runner into the given directory.
// The cloned repo uses the same homeDir and binDir as its origin.
func Clone(original testshell.Mocking, targetDir string) (TestRuntime, error) {
	_, err := original.Run("git", "clone", original.WorkingDir, targetDir)
	if err != nil {
		return TestRuntime{}, fmt.Errorf("cannot clone repo %q: %w", original.WorkingDir, err)
	}
	return New(targetDir, original.HomeDir, original.BinDir), nil
}
