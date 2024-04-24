package format_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/format"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestStringerPtrSettingTest(t *testing.T) {
	t.Parallel()
	t.Run("valid value", func(t *testing.T) {
		t.Parallel()
		give := configdomain.GitHubToken("token")
		have := format.StringerPtrSetting(&give)
		want := "token"
		must.EqOp(t, want, have)
	})
	t.Run("direct nil value", func(t *testing.T) {
		t.Parallel()
		have := format.StringerPtrSetting(nil)
		want := "(not set)"
		must.EqOp(t, want, have)
	})
	t.Run("indirect nil value", func(t *testing.T) {
		t.Parallel()
		var give *configdomain.GitHubToken
		have := format.StringerPtrSetting(give)
		want := "(not set)"
		must.EqOp(t, want, have)
	})
}
