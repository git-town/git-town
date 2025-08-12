package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestSingleSnapshot(t *testing.T) {

	t.Run("DefaultBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.SetDefaultGitBranch("main")
		have := gitconfig.DefaultBranch(runtime)
		want := gitdomain.NewLocalBranchNameOption("main")
		must.Eq(t, want, have)
	})
}
