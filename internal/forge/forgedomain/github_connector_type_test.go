package forgedomain_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseGitHubConnectorType(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		val Option[forgedomain.GitHubConnectorType]
		err error
	}{
		"api": {
			val: Some(forgedomain.GitHubConnectorTypeAPI),
			err: nil,
		},
		"gh": {
			val: Some(forgedomain.GitHubConnectorTypeGh),
			err: nil,
		},
		"invalid": {
			val: None[forgedomain.GitHubConnectorType](),
			err: errors.New(`unknown GitHubConnectorType: "invalid"`),
		},
		"": {
			val: None[forgedomain.GitHubConnectorType](),
			err: nil,
		},
	}
	for give, want := range tests {
		have, err := forgedomain.ParseGitHubConnectorType(give)
		must.Eq(t, want.err, err)
		must.Eq(t, want.val, have)
	}
}
