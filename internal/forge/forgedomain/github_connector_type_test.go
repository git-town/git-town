package forgedomain_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseGitHubConnectorType(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		val Option[forgedomain.GithubConnectorType]
		err error
	}{
		"api": {
			val: Some(forgedomain.GithubConnectorTypeAPI),
			err: nil,
		},
		"gh": {
			val: Some(forgedomain.GithubConnectorTypeGh),
			err: nil,
		},
		"invalid": {
			val: None[forgedomain.GithubConnectorType](),
			err: errors.New(`unknown GitHubConnectorType defined in test: "invalid"`),
		},
		"": {
			val: None[forgedomain.GithubConnectorType](),
			err: nil,
		},
	}
	for give, want := range tests {
		have, err := forgedomain.ParseGithubConnectorType(give, "test")
		must.Eq(t, want.err, err)
		must.Eq(t, want.val, have)
	}
}
