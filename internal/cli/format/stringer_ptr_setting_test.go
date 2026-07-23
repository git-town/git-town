package format_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/cli/format"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/shoenig/test/must"
)

func TestOptionalStringerSetting(t *testing.T) {
	t.Parallel()
	tests := map[stringss.Trimmed]string{
		"my token": "my token",
		"":         "(not set)",
	}
	for give, want := range tests {
		option := forgedomain.ParseGithubToken(give)
		have := format.OptionalStringerSetting(option)
		must.EqOp(t, want, have)
	}
}
