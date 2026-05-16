package format_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/cli/format"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestOptionalStringerSetting(t *testing.T) {
	t.Parallel()
	tests := map[Option[string]]string{
		Some("my token"): "my token",
		None[string]():   "(not set)",
	}
	for give, want := range tests {
		option := forgedomain.ParseGithubToken(give)
		have := format.OptionalStringerSetting(option)
		must.EqOp(t, want, have)
	}
}
