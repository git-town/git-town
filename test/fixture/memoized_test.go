package fixture_test

import (
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v19/test/asserts"
	"github.com/git-town/git-town/v19/test/fixture"
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
