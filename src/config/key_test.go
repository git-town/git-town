package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/shoenig/test"
)

func TestKey(t *testing.T) {
	t.Parallel()

	t.Run("ParseKey", func(t *testing.T) {
		t.Parallel()
		t.Run("normal config key", func(t *testing.T) {
			t.Parallel()
			have := config.ParseKey("git-town.offline")
			want := &config.KeyOffline
			test.EqOp(t, want, have)
		})
		t.Run("lineage keys", func(t *testing.T) {
			t.Parallel()
			t.Run("valid lineage key", func(t *testing.T) {
				t.Parallel()
				give := "git-town-branch.branch-1.parent"
				have := config.ParseKey(give)
				want := &config.Key{give}
				test.EqOp(t, want, have)
			})
			t.Run("lineage key without suffix", func(t *testing.T) {
				t.Parallel()
				have := config.ParseKey("git-town-branch.branch-1")
				test.Nil(t, have)
			})
			t.Run("lineage key without prefix", func(t *testing.T) {
				t.Parallel()
				have := config.ParseKey("git-town.branch-1.parent")
				test.Nil(t, have)
			})
		})
		t.Run("alias key", func(t *testing.T) {
			t.Parallel()
			t.Run("valid alias", func(t *testing.T) {
				t.Parallel()
				have := config.ParseKey("alias.append")
				want := &config.KeyAliasAppend
				test.EqOp(t, want, have)
			})
			t.Run("invalid alias", func(t *testing.T) {
				t.Parallel()
				have := config.ParseKey("alias.zonk")
				test.Nil(t, have)
			})
		})
		t.Run("unknown key", func(t *testing.T) {
			t.Parallel()
			have := config.ParseKey("zonk")
			test.Nil(t, have)
		})
	})
}
