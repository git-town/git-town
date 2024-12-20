package fixture_test

import (
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/test/asserts"
	"github.com/git-town/git-town/v17/test/fixture"
	"github.com/shoenig/test/must"
)

func TestNewFreshFixture(t *testing.T) {
	t.Parallel()
	gitEnvRootDir := t.TempDir()
	result := fixture.NewFresh(gitEnvRootDir).AsFixture()
	devRepo := result.DevRepo.GetOrPanic()
	// verify the developer repo
	asserts.IsGitRepo(t, filepath.Join(gitEnvRootDir, "developer"))
	branch, err := devRepo.CurrentBranch(devRepo.TestRunner)
	must.NoError(t, err)
	must.EqOp(t, gitdomain.NewLocalBranchName("initial"), branch)
}
