package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestAliasableCommand(t *testing.T) {
	t.Parallel()

	t.Run("NewAliasableCommand", func(t *testing.T) {
		t.Parallel()
		give := "append"
		have := configdomain.NewAliasableCommand(give)
		want := configdomain.AliasableCommandAppend
		must.Eq(t, want, have)
	})

	t.Run("NewAliasableCommands", func(t *testing.T) {
		t.Parallel()
		give := []string{"append", "diff-parent"}
		have := configdomain.NewAliasableCommands(give...)
		want := configdomain.AliasableCommands{
			configdomain.AliasableCommandAppend,
			configdomain.AliasableCommandDiffParent,
		}
		must.Eq(t, want, have)
	})

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
