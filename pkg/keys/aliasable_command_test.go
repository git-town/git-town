package keys_test

import (
	"testing"

	"github.com/git-town/git-town/v14/pkg/keys"
	"github.com/shoenig/test/must"
)

func TestAliasableCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key", func(t *testing.T) {
		t.Parallel()
		must.EqOp(t, keys.AliasableCommandAppend.Key().Key(), keys.KeyAliasAppend)
		must.EqOp(t, keys.AliasableCommandSync.Key().Key(), keys.KeyAliasSync)
	})

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		give := keys.AliasableCommands{
			keys.AliasableCommandAppend,
			keys.AliasableCommandDiffParent,
			keys.AliasableCommandHack,
		}
		have := give.Strings()
		want := []string{"append", "diff-parent", "hack"}
		must.Eq(t, want, have)
	})
}