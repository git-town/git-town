package testruntime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/cache"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/test/commands"
	"github.com/git-town/git-town/v21/internal/test/subshell"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

// TestRuntime provides Git functionality for test code (unit and end-to-end tests).
type TestRuntime struct {
	commands.TestCommands
}

// Clone creates a clone of the repository managed by this test.Runner into the given directory.
// The cloned repo uses the same homeDir and binDir as its origin.
func Clone(original *subshell.TestRunner, targetDir string) commands.TestCommands {
	original.MustRun("git", "clone", original.WorkingDir, targetDir)
	return New(targetDir, original.HomeDir, original.BinDir)
}

// Create creates test.Runner instances.
func Create(t *testing.T) commands.TestCommands {
	t.Helper()
	dir := t.TempDir()
	workingDir := filepath.Join(dir, "repo")
	err := os.Mkdir(workingDir, 0o744)
	must.NoError(t, err)
	homeDir := filepath.Join(dir, "home")
	err = os.Mkdir(homeDir, 0o744)
	must.NoError(t, err)
	runtime := Initialize(workingDir, homeDir, homeDir)
	must.NoError(t, err)
	return runtime
}

// CreateGitTown creates a test.Runtime for use in tests,
// with a main branch and initial git town configuration.
func CreateGitTown(t *testing.T) commands.TestCommands {
	t.Helper()
	repo := Create(t)
	repo.CreateBranch("main", "initial")
	err := repo.Config.SetMainBranch("main", repo.TestRunner)
	must.NoError(t, err)
	err = gitconfig.SetPerennialBranches(repo.TestRunner, gitdomain.LocalBranchNames{}, configdomain.ConfigScopeLocal)
	must.NoError(t, err)
	return repo
}

// Initialize creates a fully functioning test.Runner in the given working directory,
// including necessary Git configuration to make commits. Creates missing folders as needed.
func Initialize(workingDir, homeDir, binDir string) commands.TestCommands {
	runtime := InitializeNoInitialCommit(workingDir, homeDir, binDir)
	runtime.MustRun("git", "commit", "--allow-empty", "-m", "initial commit")
	return runtime
}

// InitializeNoInitialCommit creates a fully functioning test.Runner in the given working directory,
// including necessary Git configuration to make commits. Creates missing folders as needed.
// Does not create an initial commit.
// This is useful for scenarios that require testing the behavior of Git Town in a fresh repository.
func InitializeNoInitialCommit(workingDir, homeDir, binDir string) commands.TestCommands {
	runtime := New(workingDir, homeDir, binDir)
	runtime.MustRun("git", "init", "--initial-branch=initial")
	runtime.MustRun("git", "config", "--global", "user.name", "user")
	runtime.MustRun("git", "config", "--global", "user.email", "email@example.com")
	return runtime
}

// New provides a new test.Runner instance working in the given directory.
// The directory must contain an existing Git repo.
func New(workingDir, homeDir, binDir string) commands.TestCommands {
	testRunner := subshell.TestRunner{
		BinDir:           binDir,
		HomeDir:          homeDir,
		ProposalOverride: None[string](),
		Verbose:          false,
		WorkingDir:       workingDir,
	}
	gitCommands := git.Commands{
		CurrentBranchCache: &cache.WithPrevious[gitdomain.LocalBranchName]{},
		RemotesCache:       &cache.Cache[gitdomain.Remotes]{},
	}
	unvalidatedConfig := config.NewUnvalidatedConfig(config.NewUnvalidatedConfigArgs{
		CliConfig:     configdomain.EmptyPartialConfig(),
		ConfigFile:    configdomain.EmptyPartialConfig(),
		Defaults:      config.DefaultNormalConfig(),
		EnvConfig:     configdomain.EmptyPartialConfig(),
		FinalMessages: stringslice.NewCollector(),
		GitGlobal:     configdomain.EmptyPartialConfig(),
		GitLocal:      configdomain.EmptyPartialConfig(),
		GitUnscoped:   configdomain.EmptyPartialConfig(),
	})
	unvalidatedConfig.UnvalidatedConfig.MainBranch = gitdomain.NewLocalBranchNameOption("main")
	testCommands := commands.TestCommands{
		Git:        &gitCommands,
		Config:     unvalidatedConfig,
		SnapShots:  map[configdomain.ConfigScope]configdomain.SingleSnapshot{},
		TestRunner: &testRunner,
	}
	return testCommands
}
