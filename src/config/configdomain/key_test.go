package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestKey(t *testing.T) {
	t.Parallel()

	t.Run("IsAliasKey", func(t *testing.T) {
		t.Parallel()
		must.True(t, configdomain.KeyAliasAppend.IsAliasKey())
		must.False(t, configdomain.KeyPushHook.IsAliasKey())
	})

	t.Run("IsLineage", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[configdomain.LineageKey]{
			"git-town-branch.branch.parent": Some(configdomain.LineageKey("git-town-branch.branch.parent")), // valid lineage key
			"git-town-branch..parent":       Some(configdomain.LineageKey("git-town-branch..parent")),       // empty lineage key
			"git-town.push-hook":            None[configdomain.LineageKey](),                                // not a lineage key
		}
		for give, want := range tests {
			key := configdomain.Key(give)
			have := key.IsLineage()
			must.Eq(t, want, have)
		}
	})

	t.Run("ParseKey", func(t *testing.T) {
		t.Parallel()
		t.Run("normal config key", func(t *testing.T) {
			t.Parallel()
			have, has := configdomain.ParseKey("git-town.offline").Get()
			must.True(t, has)
			want := configdomain.KeyOffline
			must.EqOp(t, want, have)
		})
		t.Run("lineage keys", func(t *testing.T) {
			t.Parallel()
			t.Run("valid lineage key", func(t *testing.T) {
				t.Parallel()
				give := "git-town-branch.branch-1.parent"
				have, has := configdomain.ParseKey(give).Get()
				must.True(t, has)
				want := configdomain.Key(give)
				must.EqOp(t, want, have)
			})
			t.Run("lineage key without suffix", func(t *testing.T) {
				t.Parallel()
				have := configdomain.ParseKey("git-town-branch.branch-1")
				must.True(t, have.IsNone())
			})
			t.Run("lineage key without prefix", func(t *testing.T) {
				t.Parallel()
				have := configdomain.ParseKey("git-town.branch-1.parent")
				must.True(t, have.IsNone())
			})
		})
		t.Run("alias key", func(t *testing.T) {
			t.Parallel()
			t.Run("valid alias", func(t *testing.T) {
				t.Parallel()
				have, has := configdomain.ParseKey("alias.append").Get()
				must.True(t, has)
				must.NotNil(t, have)
				want := configdomain.KeyAliasAppend
				must.EqOp(t, want, have)
			})
			t.Run("invalid alias", func(t *testing.T) {
				t.Parallel()
				have := configdomain.ParseKey("alias.zonk")
				must.True(t, have.IsNone())
			})
		})
		t.Run("unknown key", func(t *testing.T) {
			t.Parallel()
			have := configdomain.ParseKey("zonk")
			must.True(t, have.IsNone())
		})
	})
}
