package format_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/cli/format"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestOptionalStringerSetting(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"my token": "my token",
		"":         "(not set)",
	}
	for give, want := range tests {
		option := configdomain.ParseGitHubToken(give)
		have := format.OptionalStringerSetting(option)
		must.EqOp(t, want, have)
	}
}
