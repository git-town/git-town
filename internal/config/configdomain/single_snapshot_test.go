package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/test/testruntime"
	"github.com/git-town/git-town/v21/pkg/asserts"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestSingleSnapshot(t *testing.T) {
	t.Run("DefaultBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.SetDefaultGitBranch("custom")
		snapshot := asserts.NoError1(gitconfig.LoadSnapshot(runtime.TestRunner, Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedNo))
		have := snapshot.DefaultBranch(runtime)
		want := gitdomain.NewLocalBranchNameOption("main")
		must.Eq(t, want, have)
	})
}
