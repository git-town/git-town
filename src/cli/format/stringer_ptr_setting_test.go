package format_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/format"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestOptionalStringerSetting(t *testing.T) {
	t.Parallel()
	tests := map[gohacks.Option[configdomain.GitHubToken]]string{
		gohacks.NewOption(configdomain.NewGitHubToken("my token")): "my token",
		gohacks.NewOptionNone[configdomain.GitHubToken]():          "(not set)",
	}
	for give, want := range tests {
		have := format.OptionalStringerSetting[configdomain.GitHubToken](give)
		must.EqOp(t, want, have)
	}
}
