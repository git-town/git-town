package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestAliasableCommand(t *testing.T) {
	t.Parallel()

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		give := configdomain.AliasableCommands{
			configdomain.AliasableCommandAppend,
			configdomain.AliasableCommandDiffParent,
			configdomain.AliasableCommandHack,
		}
		have := give.Strings()
		want := []string{"append", "diff-parent", "hack"}
		must.Eq(t, want, have)
	})
}
