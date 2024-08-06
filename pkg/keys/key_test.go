package keys_test

import (
	"testing"

	"github.com/git-town/git-town/v14/pkg/keys"
	"github.com/shoenig/test/must"
)

func TestKey(t *testing.T) {
	t.Parallel()

	t.Run("ParseKey", func(t *testing.T) {
		t.Parallel()
		t.Run("normal config key", func(t *testing.T) {
			t.Parallel()
			have, has := keys.ParseKey("git-town.offline").Get()
			must.True(t, has)
			want := keys.KeyOffline
			must.EqOp(t, want, have)
		})
		t.Run("lineage keys", func(t *testing.T) {
			t.Parallel()
			t.Run("valid lineage key", func(t *testing.T) {
				t.Parallel()
				give := "git-town-branch.branch-1.parent"
				have, has := keys.ParseKey(give).Get()
				must.True(t, has)
				want := keys.Key(give)
				must.EqOp(t, want, have)
			})
			t.Run("lineage key without suffix", func(t *testing.T) {
				t.Parallel()
				have := keys.ParseKey("git-town-branch.branch-1")
				must.True(t, have.IsNone())
			})
			t.Run("lineage key without prefix", func(t *testing.T) {
				t.Parallel()
				have := keys.ParseKey("git-town.branch-1.parent")
				must.True(t, have.IsNone())
			})
		})
		t.Run("alias key", func(t *testing.T) {
			t.Parallel()
			t.Run("valid alias", func(t *testing.T) {
				t.Parallel()
				have, has := keys.ParseKey("alias.append").Get()
				must.True(t, has)
				must.NotNil(t, have)
				want := keys.KeyAliasAppend
				must.EqOp(t, want, have)
			})
			t.Run("invalid alias", func(t *testing.T) {
				t.Parallel()
				have := keys.ParseKey("alias.zonk")
				must.True(t, have.IsNone())
			})
		})
		t.Run("unknown key", func(t *testing.T) {
			t.Parallel()
			have := keys.ParseKey("zonk")
			must.True(t, have.IsNone())
		})
	})
}
