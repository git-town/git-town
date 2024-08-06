package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestAliasKey(t *testing.T) {
	t.Parallel()

	t.Run("NewAliasKey", func(t *testing.T) {
		t.Parallel()
		tests := map[configdomain.Key]Option[configdomain.AliasKey]{
			configdomain.KeyAliasAppend: Some(configdomain.AliasKey(configdomain.KeyAliasAppend)),
			configdomain.KeyPushHook:    None[configdomain.AliasKey](),
		}
		for give, want := range tests {
			have := configdomain.NewAliasKey(give)
			must.Eq(t, want, have)
		}
	})
}
