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
		t.Run("nothing changed", func(t *testing.T) {
			t.Parallel()
			have := configdomain.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeyGithubToken,
				},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			other := configdomain.ConfigDiff{ //nolint:exhaustruct
			}
			have.Merge(&other)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("added entries", func(t *testing.T) {
			t.Parallel()
			have := configdomain.ConfigDiff{
				Added: []configdomain.Key{
					configdomain.KeyGiteaToken,
				},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			other := configdomain.ConfigDiff{ //nolint:exhaustruct
				Added: []configdomain.Key{
					configdomain.KeyGithubToken,
				},
			}
			have.Merge(&other)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGiteaToken, configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed entries", func(t *testing.T) {
			t.Parallel()
			have := configdomain.ConfigDiff{
				Removed: map[configdomain.Key]string{
					configdomain.KeyGiteaToken: "gitea",
				},
				Added:   []configdomain.Key{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			other := configdomain.ConfigDiff{ //nolint:exhaustruct
				Removed: map[configdomain.Key]string{
					configdomain.KeyGithubToken: "github",
				},
			}
			have.Merge(&other)
			want := configdomain.ConfigDiff{
				Removed: map[configdomain.Key]string{
					configdomain.KeyGiteaToken:  "gitea",
					configdomain.KeyGithubToken: "github",
				},
				Added:   []configdomain.Key{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed entries", func(t *testing.T) {
			t.Parallel()
			have := configdomain.ConfigDiff{
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyGiteaToken: {
						Before: "giteaBefore",
						After:  "giteaAfter",
					},
				},
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
			}
			other := configdomain.ConfigDiff{ //nolint:exhaustruct
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyGithubToken: {
						Before: "githubBefore",
						After:  "githubAfter",
					},
				},
			}
			have.Merge(&other)
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
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("CheckPtr", func(t *testing.T) {
		t.Parallel()
		t.Run("value not changed", func(t *testing.T) {
			t.Parallel()
			before := domain.LocalBranchName("main")
			after := domain.LocalBranchName("main")
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("added", func(t *testing.T) {
			t.Parallel()
			after := configdomain.GitHubToken("token")
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffPtr(&have, configdomain.KeyGithubToken, nil, &after)
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
			configdomain.DiffPtr(&have, configdomain.KeyGithubToken, &before, nil)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{configdomain.KeyGithubToken: "token"},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed", func(t *testing.T) {
			t.Parallel()
			before := configdomain.GitHubToken("token1")
			after := configdomain.GitHubToken("token2")
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyGithubToken: {
						Before: "token1",
						After:  "token2",
					},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed from empty string", func(t *testing.T) {
			t.Parallel()
			before := configdomain.GitHubToken("")
			after := configdomain.GitHubToken("token")
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyGithubToken: {
						Before: "",
						After:  "token",
					},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed to empty string", func(t *testing.T) {
			t.Parallel()
			before := configdomain.GitHubToken("token")
			after := configdomain.GitHubToken("")
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyGithubToken: {
						Before: "token",
						After:  "",
					},
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("CheckString", func(t *testing.T) {
		t.Parallel()
		t.Run("value not changed", func(t *testing.T) {
			t.Parallel()
			before := "main"
			after := "main"
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffString(&have, configdomain.KeyGithubToken, before, after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("added", func(t *testing.T) {
			t.Parallel()
			before := ""
			after := "token"
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffString(&have, configdomain.KeyGithubToken, before, after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed", func(t *testing.T) {
			t.Parallel()
			before := "token"
			after := ""
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffString(&have, configdomain.KeyGithubToken, before, after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{configdomain.KeyGithubToken: "token"},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed", func(t *testing.T) {
			t.Parallel()
			before := "token1"
			after := "token2"
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffString(&have, configdomain.KeyGithubToken, before, after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyGithubToken: {
						Before: "token1",
						After:  "token2",
					},
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("CheckStringPtr", func(t *testing.T) {
		t.Parallel()
		t.Run("value not changed", func(t *testing.T) {
			t.Parallel()
			before := "main"
			after := "main"
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffStringPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("added to nil", func(t *testing.T) {
			t.Parallel()
			after := "token"
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffStringPtr(&have, configdomain.KeyGithubToken, nil, &after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("added to empty string", func(t *testing.T) {
			t.Parallel()
			before := ""
			after := "token"
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffStringPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed to nil", func(t *testing.T) {
			t.Parallel()
			before := "token"
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffStringPtr(&have, configdomain.KeyGithubToken, &before, nil)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{configdomain.KeyGithubToken: "token"},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed to empty string", func(t *testing.T) {
			t.Parallel()
			before := "token"
			after := ""
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffStringPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{configdomain.KeyGithubToken: "token"},
				Changed: map[configdomain.Key]domain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed", func(t *testing.T) {
			t.Parallel()
			before := "token1"
			after := "token2"
			have := configdomain.EmptyConfigDiff()
			configdomain.DiffStringPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := configdomain.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]domain.Change[string]{
					configdomain.KeyGithubToken: {
						Before: "token1",
						After:  "token2",
					},
				},
			}
			must.Eq(t, want, have)
		})
	})
}
