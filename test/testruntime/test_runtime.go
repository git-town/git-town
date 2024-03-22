package testruntime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v12/src/config"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/gohacks/cache"
	"github.com/git-town/git-town/v12/test/commands"
	testshell "github.com/git-town/git-town/v12/test/subshell"
	"github.com/shoenig/test/must"
)

// TestRuntime provides Git functionality for test code (unit and end-to-end tests).
type TestRuntime struct {
	commands.TestCommands
	Backend git.BackendCommands
}

// Clone creates a clone of the repository managed by this test.Runner into the given directory.
// The cloned repo uses the same homeDir and binDir as its origin.
func Clone(original *testshell.TestRunner, targetDir string) TestRuntime {
	original.MustRun("git", "clone", original.WorkingDir, targetDir)
	return New(targetDir, original.HomeDir, original.BinDir)
}

// Create creates test.Runner instances.
func Create(t *testing.T) TestRuntime {
	t.Helper()
	dir := t.TempDir()
	workingDir := filepath.Join(dir, "repo")
	err := os.Mkdir(workingDir, 0o744)
	must.NoError(t, err)
	homeDir := filepath.Join(dir, "home")
	err = os.Mkdir(homeDir, 0o744)
	must.NoError(t, err)
	runtime := Initialize(workingDir, homeDir, homeDir)
	err = runtime.Run("git", "commit", "--allow-empty", "-m", "initial commit")
	must.NoError(t, err)
	return runtime
}

// CreateGitTown creates a test.Runtime for use in tests,
// with a main branch and initial git town configuration.
func CreateGitTown(t *testing.T) TestRuntime {
	t.Helper()
	repo := Create(t)
	repo.CreateBranch(gitdomain.NewLocalBranchName("main"), gitdomain.NewLocalBranchName("initial"))
	err := repo.Config.SetMainBranch(gitdomain.NewLocalBranchName("main"))
	must.NoError(t, err)
	err = repo.Config.SetPerennialBranches(gitdomain.LocalBranchNames{})
	must.NoError(t, err)
	return repo
}

// initialize creates a fully functioning test.Runner in the given working directory,
// including necessary Git configuration to make commits. Creates missing folders as needed.
func Initialize(workingDir, homeDir, binDir string) TestRuntime {
	runtime := New(workingDir, homeDir, binDir)
	runtime.MustRunMany([][]string{
		{"git", "init", "--initial-branch=initial"},
		{"git", "config", "--global", "user.name", "user"},
		{"git", "config", "--global", "user.email", "email@example.com"},
	})
	return runtime
}

// newRuntime provides a new test.Runner instance working in the given directory.
// The directory must contain an existing Git repo.
func New(workingDir, homeDir, binDir string) TestRuntime {
	runner := testshell.TestRunner{
		WorkingDir: workingDir,
		HomeDir:    homeDir,
		BinDir:     binDir,
	}
	config, _, err := config.NewConfig(config.NewConfigArgs{
		ConfigFile:   nil,
		DryRun:       false,
		GlobalConfig: configdomain.EmptyPartialConfig(),
		LocalConfig:  configdomain.EmptyPartialConfig(),
		Runner:       &runner,
	})
	if err != nil {
		panic(err)
	}
	backendCommands := git.BackendCommands{
		Runner:             &runner,
		DryRun:             false,
		Config:             config,
		CurrentBranchCache: &cache.LocalBranchWithPrevious{},
		RemotesCache:       &cache.Remotes{},
	}
	testCommands := commands.TestCommands{
		BackendCommands: &backendCommands,
		TestRunner:      &runner,
	}
	return TestRuntime{
		Backend:      backendCommands,
		TestCommands: testCommands,
	}
}
