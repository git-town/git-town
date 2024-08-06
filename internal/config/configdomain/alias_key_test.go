package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/pkg/keys"
	. "github.com/git-town/git-town/v14/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestAliasKey(t *testing.T) {

	t.Run("NewAliasKey", func(t *testing.T) {
		t.Parallel()
		tests := map[keys.Key]Option[configdomain.AliasKey]{
			keys.KeyAliasAppend: Some(configdomain.AliasKey(keys.KeyAliasAppend)),
			keys.KeyPushHook:    None[configdomain.AliasKey](),
		}
		for give, want := range tests {
			have := configdomain.NewAliasKey(give)
			must.Eq(t, want, have)
		}
	})
}
