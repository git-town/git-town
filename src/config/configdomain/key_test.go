package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestKey(t *testing.T) {
	t.Parallel()

	t.Run("ParseKey", func(t *testing.T) {
		t.Parallel()
		t.Run("normal config key", func(t *testing.T) {
			t.Parallel()
			have := configdomain.ParseKey("git-town.offline")
			want := &configdomain.KeyOffline
			must.EqOp(t, *want, *have)
		})
		t.Run("lineage keys", func(t *testing.T) {
			t.Parallel()
			t.Run("valid lineage key", func(t *testing.T) {
				t.Parallel()
				give := "git-town-branch.branch-1.parent"
				have := configdomain.ParseKey(give)
				want := configdomain.NewKey(give)
				must.EqOp(t, want, *have)
			})
			t.Run("lineage key without suffix", func(t *testing.T) {
				t.Parallel()
				have := configdomain.ParseKey("git-town-branch.branch-1")
				must.Nil(t, have)
			})
			t.Run("lineage key without prefix", func(t *testing.T) {
				t.Parallel()
				have := configdomain.ParseKey("git-town.branch-1.parent")
				must.Nil(t, have)
			})
		})
		t.Run("alias key", func(t *testing.T) {
			t.Parallel()
			t.Run("valid alias", func(t *testing.T) {
				t.Parallel()
				have := configdomain.ParseKey("alias.append")
				want := &configdomain.KeyAliasAppend
				must.EqOp(t, *want, *have)
			})
			t.Run("invalid alias", func(t *testing.T) {
				t.Parallel()
				have := configdomain.ParseKey("alias.zonk")
				must.Nil(t, have)
			})
		})
		t.Run("unknown key", func(t *testing.T) {
			t.Parallel()
			have := configdomain.ParseKey("zonk")
			must.Nil(t, have)
		})
	})
}
