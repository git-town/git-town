package fixture_test

import (
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v22/internal/test/fixture"
	"github.com/git-town/git-town/v22/pkg/asserts"
	"github.com/shoenig/test/must"
)

func TestNewFreshFixture(t *testing.T) {
	t.Parallel()
	gitEnvRootDir := t.TempDir()
	result := fixture.NewFresh(gitEnvRootDir).AsFixture()
	devRepo := result.DevRepo.GetOrPanic()
	// verify the developer repo
	asserts.IsGitRepo(t, filepath.Join(gitEnvRootDir, "developer"))
	branch := asserts.NoError1(devRepo.Git.CurrentBranch(devRepo.TestRunner)).GetOrPanic()
	must.EqOp(t, "initial", branch)
}
