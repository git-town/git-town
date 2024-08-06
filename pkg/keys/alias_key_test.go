package keys_test

import (
	"testing"

	"github.com/git-town/git-town/v14/pkg/keys"
	. "github.com/git-town/git-town/v14/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestAliasKey(t *testing.T) {

	t.Run("NewAliasKey", func(t *testing.T) {
		t.Parallel()
		tests := map[keys.Key]Option[keys.AliasKey]{
			keys.KeyAliasAppend: Some(keys.AliasKey(keys.KeyAliasAppend)),
			keys.KeyPushHook:    None[keys.AliasKey](),
		}
		for give, want := range tests {
			have := keys.NewAliasKey(give)
			must.Eq(t, want, have)
		}
	})
}
