package testruntime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v15/internal/config"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/config/gitconfig"
	"github.com/git-town/git-town/v15/internal/git"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/gohacks/cache"
	. "github.com/git-town/git-town/v15/pkg/prelude"
	"github.com/git-town/git-town/v15/test/commands"
	testshell "github.com/git-town/git-town/v15/test/subshell"
	"github.com/shoenig/test/must"
)

// TestRuntime provides Git functionality for test code (unit and end-to-end tests).
type TestRuntime struct {
	commands.TestCommands
	Config config.ValidatedConfig
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
	must.NoError(t, err)
	return runtime
}

// CreateGitTown creates a test.Runtime for use in tests,
// with a main branch and initial git town configuration.
func CreateGitTown(t *testing.T) TestRuntime {
	t.Helper()
	repo := Create(t)
	repo.CreateBranch(gitdomain.NewLocalBranchName("main"), gitdomain.NewBranchName("initial"))
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
	runtime.MustRun("git", "init", "--initial-branch=initial")
	runtime.MustRun("git", "config", "--global", "user.name", "user")
	runtime.MustRun("git", "config", "--global", "user.email", "email@example.com")
	runtime.MustRun("git", "commit", "--allow-empty", "-m", "initial commit")
	return runtime
}

// newRuntime provides a new test.Runner instance working in the given directory.
// The directory must contain an existing Git repo.
func New(workingDir, homeDir, binDir string) TestRuntime {
	testRunner := testshell.TestRunner{
		BinDir:           binDir,
		HomeDir:          homeDir,
		ProposalOverride: None[string](),
		Verbose:          false,
		WorkingDir:       workingDir,
	}
	gitCommands := git.Commands{
		CurrentBranchCache: &cache.LocalBranchWithPrevious{},
		RemotesCache:       &cache.Remotes{},
	}
	unvalidatedConfig, _ := config.NewUnvalidatedConfig(config.NewUnvalidatedConfigArgs{
		Access: gitconfig.Access{
			Runner: &testRunner,
		},
		ConfigFile:   None[configdomain.PartialConfig](),
		DryRun:       false,
		GlobalConfig: configdomain.EmptyPartialConfig(),
		LocalConfig:  configdomain.EmptyPartialConfig(),
	})
	validatedConfig := config.ValidatedConfig{
		Config: configdomain.ValidatedConfig{
			UnvalidatedConfig: unvalidatedConfig.Config.Value,
			GitUserEmail:      "test@test.com",
			GitUserName:       "Tester",
			MainBranch:        gitdomain.NewLocalBranchName("main"),
		},
		UnvalidatedConfig: &unvalidatedConfig,
	}
	testCommands := commands.TestCommands{
		Commands:   &gitCommands,
		Config:     validatedConfig,
		TestRunner: &testRunner,
	}
	return TestRuntime{
		Config:       validatedConfig,
		TestCommands: testCommands,
	}
}
