package git_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/git"
	"github.com/shoenig/test/must"
)

func TestIsAcceptableGitVersion(t *testing.T) {
	t.Parallel()
	tests := map[git.Version]bool{
		{2, 30}: true,
		{3, 0}:  true,
		{2, 29}: false,
		{1, 8}:  false,
	}
	for version, want := range tests {
		have := version.IsMinimumRequiredGitVersion()
		must.EqOp(t, want, have)
	}
}
