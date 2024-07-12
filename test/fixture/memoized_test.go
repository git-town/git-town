package fixture_test

import (
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/test/asserts"
	"github.com/git-town/git-town/v14/test/fixture"
	"github.com/shoenig/test/must"
)

func TestNewStandardFixture(t *testing.T) {
	t.Parallel()
	gitEnvRootDir := t.TempDir()
	result := fixture.NewMemoized(gitEnvRootDir).AsFixture()
	// verify the origin repo
	asserts.IsGitRepo(t, filepath.Join(gitEnvRootDir, "origin"))
	branch, err := result.OriginRepo.GetOrPanic().CurrentBranch(result.DevRepo.TestRunner)
	must.NoError(t, err)
	must.EqOp(t, gitdomain.NewLocalBranchName("main"), branch)
	// verify the developer repo
	asserts.IsGitRepo(t, filepath.Join(gitEnvRootDir, "developer"))
	assertHasGitConfiguration(t, gitEnvRootDir)
	branch, err = result.DevRepo.CurrentBranch(result.DevRepo.TestRunner)
	must.NoError(t, err)
	must.EqOp(t, gitdomain.NewLocalBranchName("main"), branch)
}
