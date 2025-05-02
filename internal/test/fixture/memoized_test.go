package fixture_test

import (
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v20/internal/test/fixture"
	"github.com/git-town/git-town/v20/pkg/asserts"
	"github.com/shoenig/test/must"
)

func TestNewStandardFixture(t *testing.T) {
	t.Parallel()
	gitEnvRootDir := t.TempDir()
	result := fixture.NewMemoized(gitEnvRootDir).AsFixture()
	devRepo := result.DevRepo.GetOrPanic()
	// verify the origin repo
	asserts.IsGitRepo(t, filepath.Join(gitEnvRootDir, "origin"))
	branch, err := result.OriginRepo.GetOrPanic().Git.CurrentBranch(devRepo.TestRunner)
	must.NoError(t, err)
	must.EqOp(t, "main", branch)
	// verify the developer repo
	asserts.IsGitRepo(t, filepath.Join(gitEnvRootDir, "developer"))
	assertHasGitConfiguration(t, gitEnvRootDir)
	branch, err = devRepo.Git.CurrentBranch(devRepo.TestRunner)
	must.NoError(t, err)
	must.EqOp(t, "main", branch)
}
