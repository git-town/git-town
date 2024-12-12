package gitconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/config/gitconfig"
	"github.com/shoenig/test/must"
)

func TestIsGitTownAlias(t *testing.T) {
	t.Parallel()
	tests := map[string]bool{
		"town append": true,
		"town sync":   true,
		"other":       false,
	}
	for give, want := range tests {
		have := gitconfig.IsGitTownAlias(give)
		must.EqOp(t, want, have)
	}
}
