package config

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestFormatGithubToken(t *testing.T) {
	t.Parallel()

	t.Run("not set", func(t *testing.T) {
		t.Parallel()
		give := None[forgedomain.GithubToken]()
		have := formatGithubToken(give)
		must.EqOp(t, "(not set)", have)
	})

	t.Run("set", func(t *testing.T) {
		t.Parallel()
		give := Some(forgedomain.GithubToken("github-token"))
		have := formatGithubToken(give)
		must.EqOp(t, "(configured)", have)
	})
}
