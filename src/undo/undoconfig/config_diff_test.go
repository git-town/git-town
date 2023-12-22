package undoconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/undo/undoconfig"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/shoenig/test/must"
)

func TestConfigdiff(t *testing.T) {
	t.Parallel()

	t.Run("Merge", func(t *testing.T) {
		t.Parallel()
		t.Run("nothing changed", func(t *testing.T) {
			t.Parallel()
			have := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Added: []configdomain.Key{
					configdomain.KeyGithubToken,
				},
			}
			other := undoconfig.ConfigDiff{ //nolint:exhaustruct
			}
			have.Merge(&other)
			want := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Added: []configdomain.Key{configdomain.KeyGithubToken},
			}
			must.Eq(t, want, have)
		})
		t.Run("added entries", func(t *testing.T) {
			t.Parallel()
			have := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Added: []configdomain.Key{
					configdomain.KeyGiteaToken,
				},
			}
			other := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Added: []configdomain.Key{
					configdomain.KeyGithubToken,
				},
			}
			have.Merge(&other)
			want := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Added: []configdomain.Key{configdomain.KeyGiteaToken, configdomain.KeyGithubToken},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed entries", func(t *testing.T) {
			t.Parallel()
			have := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Removed: map[configdomain.Key]string{
					configdomain.KeyGiteaToken: "gitea",
				},
			}
			other := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Removed: map[configdomain.Key]string{
					configdomain.KeyGithubToken: "github",
				},
			}
			have.Merge(&other)
			want := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Removed: map[configdomain.Key]string{
					configdomain.KeyGiteaToken:  "gitea",
					configdomain.KeyGithubToken: "github",
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed entries", func(t *testing.T) {
			t.Parallel()
			have := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyGiteaToken: {
						Before: "giteaBefore",
						After:  "giteaAfter",
					},
				},
			}
			other := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyGithubToken: {
						Before: "githubBefore",
						After:  "githubAfter",
					},
				},
			}
			have.Merge(&other)
			want := undoconfig.ConfigDiff{ //nolint:exhaustruct
				Changed: map[configdomain.Key]undodomain.Change[string]{
					configdomain.KeyGiteaToken: {
						Before: "giteaBefore",
						After:  "giteaAfter",
					},
					configdomain.KeyGithubToken: {
						Before: "githubBefore",
						After:  "githubAfter",
					},
				},
			}
			must.Eq(t, want, have)
		})
	})
}
