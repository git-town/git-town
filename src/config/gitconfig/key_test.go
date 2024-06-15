package gitconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/shoenig/test/must"
)

func TestKey(t *testing.T) {
	t.Parallel()

	t.Run("ParseKey", func(t *testing.T) {
		t.Parallel()
		t.Run("normal config key", func(t *testing.T) {
			t.Parallel()
			have, has := gitconfig.ParseKey("git-town.offline").Get()
			must.True(t, has)
			want := gitconfig.KeyOffline
			must.EqOp(t, want, have)
		})
		t.Run("lineage keys", func(t *testing.T) {
			t.Parallel()
			t.Run("valid lineage key", func(t *testing.T) {
				t.Parallel()
				give := "git-town-branch.branch-1.parent"
				have, has := gitconfig.ParseKey(give).Get()
				must.True(t, has)
				want := gitconfig.Key(give)
				must.EqOp(t, want, have)
			})
			t.Run("lineage key without suffix", func(t *testing.T) {
				t.Parallel()
				have := gitconfig.ParseKey("git-town-branch.branch-1")
				must.Nil(t, have)
			})
			t.Run("lineage key without prefix", func(t *testing.T) {
				t.Parallel()
				have := gitconfig.ParseKey("git-town.branch-1.parent")
				must.Nil(t, have)
			})
		})
		t.Run("alias key", func(t *testing.T) {
			t.Parallel()
			t.Run("valid alias", func(t *testing.T) {
				t.Parallel()
				have, has := gitconfig.ParseKey("alias.append").Get()
				must.True(t, has)
				must.NotNil(t, have)
				want := gitconfig.KeyAliasAppend
				must.NotNil(t, want)
				must.EqOp(t, want, have)
			})
			t.Run("invalid alias", func(t *testing.T) {
				t.Parallel()
				have := gitconfig.ParseKey("alias.zonk")
				must.Nil(t, have)
			})
		})
		t.Run("unknown key", func(t *testing.T) {
			t.Parallel()
			have := gitconfig.ParseKey("zonk")
			must.Nil(t, have)
		})
	})
}
