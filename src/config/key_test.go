package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	t.Parallel()
	t.Run("ParseKey", func(t *testing.T) {
		t.Parallel()
		t.Run("normal config key", func(t *testing.T) {
			t.Parallel()
			have := config.ParseKey("git-town.offline")
			want := &config.KeyOffline
			assert.Equal(t, want, have)
		})
		t.Run("lineage keys", func(t *testing.T) {
			t.Parallel()
			t.Run("valid lineage key", func(t *testing.T) {
				t.Parallel()
				give := "git-town-branch.branch-1.parent"
				have := config.ParseKey(give)
				want := &config.Key{give}
				assert.Equal(t, want, have)
			})
			t.Run("lineage key without suffix", func(t *testing.T) {
				t.Parallel()
				have := config.ParseKey("git-town-branch.branch-1")
				assert.Nil(t, have)
			})
			t.Run("lineage key without prefix", func(t *testing.T) {
				t.Parallel()
				have := config.ParseKey("git-town.branch-1.parent")
				assert.Nil(t, have)
			})
		})
		t.Run("alias key", func(t *testing.T) {
			t.Parallel()
			t.Run("valid alias", func(t *testing.T) {
				t.Parallel()
				have := config.ParseKey("alias.append")
				want := &config.KeyAliasAppend
				assert.Equal(t, want, have)
			})
			t.Run("invalid alias", func(t *testing.T) {
				t.Parallel()
				have := config.ParseKey("alias.zonk")
				assert.Nil(t, have)
			})
		})
		t.Run("unknown key", func(t *testing.T) {
			t.Parallel()
			have := config.ParseKey("zonk")
			assert.Nil(t, have)
		})
	})
}
