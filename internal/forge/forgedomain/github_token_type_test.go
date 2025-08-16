package forgedomain_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseGitHubTokenType(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		val Option[forgedomain.GitHubTokenType]
		err error
	}{
		"enter": {
			val: Some(forgedomain.GitHubTokenTypeEnter),
			err: nil,
		},
		"cli": {
			val: Some(forgedomain.GitHubTokenTypeScript),
			err: nil,
		},
		"invalid": {
			val: None[forgedomain.GitHubTokenType](),
			err: errors.New(`unknown GitHubTokenType: "invalid"`),
		},
		"": {
			val: None[forgedomain.GitHubTokenType](),
			err: nil,
		},
	}
	for give, want := range tests {
		have, err := forgedomain.ParseGitHubTokenType(give)
		must.Eq(t, want.err, err)
		must.Eq(t, want.val, have)
	}
}
