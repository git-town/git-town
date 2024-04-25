package format_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v14/src/cli/format"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestStringerPtrSettingTest(t *testing.T) {
	t.Parallel()
	var nilInterface *configdomain.GitHubToken
	tests := map[fmt.Stringer]string{
		configdomain.NewGitHubTokenRef("token"): "token",
		nil:                                     "(not set)",
		nilInterface:                            "(not set)",
	}
	for give, want := range tests {
		have := format.StringerPtrSetting(give)
		must.EqOp(t, want, have)
	}
}
