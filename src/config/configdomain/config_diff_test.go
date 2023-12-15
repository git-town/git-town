package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestConfigdiff(t *testing.T) {
	t.Parallel()
	t.Run("Merge", func(t *testing.T) {
		t.Parallel()
		t.Run("added entries", func(t *testing.T) {
			t.Parallel()
			diff1 := configdomain.ConfigDiff{ //nolint:exhaustruct
				Added: []configdomain.Key{
					configdomain.KeyGiteaToken,
				},
			}
			diff2 := configdomain.ConfigDiff{ //nolint:exhaustruct
				Added: []configdomain.Key{
					configdomain.KeyGithubToken,
				},
			}
			have := diff1.Merge(&diff2)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGiteaToken, configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed entries", func(t *testing.T) {
			t.Parallel()
			diff1 := configdomain.ConfigDiff{ //nolint:exhaustruct
				Removed: map[configdomain.Key]string{
					configdomain.KeyGiteaToken: "gitea",
				},
			}
			diff2 := configdomain.ConfigDiff{ //nolint:exhaustruct
				Removed: map[configdomain.Key]string{
					configdomain.KeyGithubToken: "github",
				},
			}
			have := diff1.Merge(&diff2)
			want := configdomain.ConfigDiff{
				Removed: map[configdomain.Key]string{
					configdomain.KeyGiteaToken:  "gitea",
					configdomain.KeyGithubToken: "github",
				},
				Added:   nil,
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed entries", func(t *testing.T) {
			t.Parallel()
			diff1 := configdomain.ConfigDiff{ //nolint:exhaustruct
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyGiteaToken: {
						Before: "giteaBefore",
						After:  "giteaAfter",
					},
				},
			}
			diff2 := configdomain.ConfigDiff{ //nolint:exhaustruct
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyGithubToken: {
						Before: "githubBefore",
						After:  "githubAfter",
					},
				},
			}
			have := diff1.Merge(&diff2)
			want := configdomain.ConfigDiff{
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyGiteaToken: {
						Before: "giteaBefore",
						After:  "giteaAfter",
					},
					configdomain.KeyGithubToken: {
						Before: "githubBefore",
						After:  "githubAfter",
					},
				},
				Added:   nil,
				Removed: map[configdomain.Key]string{},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("Check", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("CheckPtr", func(t *testing.T) {
		t.Parallel()
		t.Run("added", func(t *testing.T) {
			t.Parallel()
			after := configdomain.GitHubToken("token")
			have := configdomain.EmptyConfigDiff()
			configdomain.CheckPtr(&have, configdomain.KeyGithubToken, nil, &after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed", func(t *testing.T) {
			t.Parallel()
			before := configdomain.GitHubToken("token")
			have := configdomain.EmptyConfigDiff()
			configdomain.CheckPtr(&have, configdomain.KeyGithubToken, &before, nil)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{configdomain.KeyGithubToken: "token"},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("CheckString", func(t *testing.T) {
		t.Parallel()
	})
}
