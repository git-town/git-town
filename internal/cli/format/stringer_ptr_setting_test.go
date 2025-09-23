package format_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/cli/format"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/shoenig/test/must"
)

func TestOptionalStringerSetting(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"my token": "my token",
		"":         "(not set)",
	}
	for give, want := range tests {
		option := forgedomain.ParseGitHubToken(give)
		have := format.OptionalStringerSetting(option)
		must.EqOp(t, want, have)
	}
}
