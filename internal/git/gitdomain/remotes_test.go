package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/test/git"
	"github.com/shoenig/test/must"
)

func TestRemotes(t *testing.T) {
	t.Parallel()

	t.Run("HasOrigin", func(t *testing.T) {
		t.Parallel()
		t.Run("origin remote exists", func(t *testing.T) {
			t.Parallel()
			remotes := gitdomain.Remotes{git.REMOTE_ORIGIN}
			must.True(t, remotes.Contains(git.REMOTE_ORIGIN))
		})
		t.Run("origin remote does not exist", func(t *testing.T) {
			t.Parallel()
			remotes := gitdomain.Remotes{gitdomain.RemoteUpstream}
			must.False(t, remotes.Contains(git.REMOTE_ORIGIN))
		})
	})

	t.Run("HasUpstream", func(t *testing.T) {
		t.Parallel()
		t.Run("upstream remote exists", func(t *testing.T) {
			t.Parallel()
			remotes := gitdomain.Remotes{gitdomain.RemoteUpstream}
			must.True(t, remotes.HasUpstream())
		})
		t.Run("upstream remote does not exist", func(t *testing.T) {
			t.Parallel()
			remotes := gitdomain.Remotes{git.REMOTE_ORIGIN}
			must.False(t, remotes.HasUpstream())
		})
	})
}
