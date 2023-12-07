package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestRemotes(t *testing.T) {
	t.Parallel()

	t.Run("HasOrigin", func(t *testing.T) {
		t.Parallel()
		t.Run("origin remote exists", func(t *testing.T) {
			t.Parallel()
			remotes := domain.Remotes{domain.OriginRemote}
			must.True(t, remotes.HasOrigin())
		})
		t.Run("origin remote does not exist", func(t *testing.T) {
			t.Parallel()
			remotes := domain.NewRemotes("foo")
			must.False(t, remotes.HasOrigin())
		})
	})

	t.Run("HasUpstream", func(t *testing.T) {
		t.Parallel()
		t.Run("upstream remote exists", func(t *testing.T) {
			t.Parallel()
			remotes := domain.Remotes{domain.UpstreamRemote}
			must.True(t, remotes.HasUpstream())
		})
		t.Run("upstream remote does not exist", func(t *testing.T) {
			t.Parallel()
			remotes := domain.NewRemotes("foo")
			must.False(t, remotes.HasUpstream())
		})
	})
}
