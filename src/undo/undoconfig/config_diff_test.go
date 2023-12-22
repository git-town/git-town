package undoconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
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

	t.Run("CheckPtr", func(t *testing.T) {
		t.Parallel()
		t.Run("value not changed", func(t *testing.T) {
			t.Parallel()
			before := gitdomain.LocalBranchName("main")
			after := gitdomain.LocalBranchName("main")
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("added", func(t *testing.T) {
			t.Parallel()
			after := configdomain.GitHubToken("token")
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffPtr(&have, configdomain.KeyGithubToken, nil, &after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed", func(t *testing.T) {
			t.Parallel()
			before := configdomain.GitHubToken("token")
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffPtr(&have, configdomain.KeyGithubToken, &before, nil)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{configdomain.KeyGithubToken: "token"},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed", func(t *testing.T) {
			t.Parallel()
			before := configdomain.GitHubToken("token1")
			after := configdomain.GitHubToken("token2")
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
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
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
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
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
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
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffString(&have, configdomain.KeyGithubToken, before, after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("added", func(t *testing.T) {
			t.Parallel()
			before := ""
			after := "token"
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffString(&have, configdomain.KeyGithubToken, before, after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed", func(t *testing.T) {
			t.Parallel()
			before := "token"
			after := ""
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffString(&have, configdomain.KeyGithubToken, before, after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{configdomain.KeyGithubToken: "token"},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed", func(t *testing.T) {
			t.Parallel()
			before := "token1"
			after := "token2"
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffString(&have, configdomain.KeyGithubToken, before, after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
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
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffStringPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("added to nil", func(t *testing.T) {
			t.Parallel()
			after := "token"
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffStringPtr(&have, configdomain.KeyGithubToken, nil, &after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("added to empty string", func(t *testing.T) {
			t.Parallel()
			before := ""
			after := "token"
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffStringPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{configdomain.KeyGithubToken},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed to nil", func(t *testing.T) {
			t.Parallel()
			before := "token"
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffStringPtr(&have, configdomain.KeyGithubToken, &before, nil)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{configdomain.KeyGithubToken: "token"},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("removed to empty string", func(t *testing.T) {
			t.Parallel()
			before := "token"
			after := ""
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffStringPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{configdomain.KeyGithubToken: "token"},
				Changed: map[configdomain.Key]undodomain.Change[string]{},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed", func(t *testing.T) {
			t.Parallel()
			before := "token1"
			after := "token2"
			have := undoconfig.EmptyConfigDiff()
			undoconfig.DiffStringPtr(&have, configdomain.KeyGithubToken, &before, &after)
			want := undoconfig.ConfigDiff{
				Added:   []configdomain.Key{},
				Removed: map[configdomain.Key]string{},
				Changed: map[configdomain.Key]undodomain.Change[string]{
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
