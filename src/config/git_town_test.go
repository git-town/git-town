package config_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/test/testruntime"
	"github.com/stretchr/testify/assert"
)

func TestGitTown(t *testing.T) {
	t.Parallel()

	t.Run("SetParent and AncestorBranches", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple ancestors", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("three", "two")
			repo.Config.SetParent("two", "one")
			repo.Config.SetParent("one", "main")
			have := repo.Config.AncestorBranches("three")
			want := []string{"main", "one", "two"}
			assert.Equal(t, want, have)
		})
		t.Run("one ancestor", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("one", "main")
			have := repo.Config.AncestorBranches("one")
			want := []string{"main"}
			assert.Equal(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("one", "main")
			have := repo.Config.AncestorBranches("two")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("BranchLineageRoots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with nested child branches", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("two", "one")
			repo.Config.SetParent("one", "main")
			repo.Config.SetParent("beta", "alpha")
			repo.Config.SetParent("alpha", "main")
			repo.Config.SetParent("hotfix1", "prod")
			repo.Config.SetParent("hotfix2", "prod")
			have := repo.Config.BranchLineageRoots()
			want := []string{"main", "prod"}
			assert.Equal(t, want, have)
		})
		t.Run("no nested branches", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("one", "main")
			repo.Config.SetParent("alpha", "main")
			have := repo.Config.BranchLineageRoots()
			want := []string{"main"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			have := repo.Config.BranchLineageRoots()
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("ChildBranches", func(t *testing.T) {
		t.Run("multiple children", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("beta1", "alpha")
			repo.Config.SetParent("beta2", "alpha")
			have := repo.Config.ChildBranches("alpha")
			want := []string{"beta1", "beta2"}
			assert.Equal(t, want, have)
		})
		t.Run("child has children", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("beta", "alpha")
			repo.Config.SetParent("gamma", "beta")
			have := repo.Config.ChildBranches("alpha")
			want := []string{"beta"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			have := repo.Config.ChildBranches("alpha")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("HasParentBranch", func(t *testing.T) {
		t.Run("has a parent", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("beta", "alpha")
			assert.True(t, repo.Config.HasParentBranch("beta"))
		})
		t.Run("has no parent", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			assert.False(t, repo.Config.HasParentBranch("foo"))
		})
	})

	t.Run("IsAncestorBranch", func(t *testing.T) {
		t.Run("greatgrandparent", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("four", "three")
			repo.Config.SetParent("three", "two")
			repo.Config.SetParent("two", "one")
			assert.True(t, repo.Config.IsAncestorBranch("four", "one"))
		})
		t.Run("direct parent", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("two", "one")
			assert.True(t, repo.Config.IsAncestorBranch("two", "one"))
		})
	})

	t.Run("ParentBranch", func(t *testing.T) {
		t.Run("has parent", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			repo.Config.SetParent("two", "one")
			assert.Equal(t, "one", repo.Config.ParentBranch("two"))
		})
		t.Run("has no parent", func(t *testing.T) {
			repo := testruntime.CreateGitTown(t)
			assert.Equal(t, "", repo.Config.ParentBranch("foo"))
		})
	})

	t.Run("OriginURL()", func(t *testing.T) {
		t.Parallel()
		tests := map[string]giturl.Parts{
			"http://github.com/organization/repository":                     {Host: "github.com", Org: "organization", Repo: "repository"},
			"http://github.com/organization/repository.git":                 {Host: "github.com", Org: "organization", Repo: "repository"},
			"https://github.com/organization/repository":                    {Host: "github.com", Org: "organization", Repo: "repository"},
			"https://github.com/organization/repository.git":                {Host: "github.com", Org: "organization", Repo: "repository"},
			"https://sub.domain.customhost.com/organization/repository":     {Host: "sub.domain.customhost.com", Org: "organization", Repo: "repository"},
			"https://sub.domain.customhost.com/organization/repository.git": {Host: "sub.domain.customhost.com", Org: "organization", Repo: "repository"},
		}
		for give, want := range tests {
			repo := testruntime.CreateGitTown(t)
			os.Setenv("GIT_TOWN_REMOTE", give)
			defer os.Unsetenv("GIT_TOWN_REMOTE")
			have := repo.Config.OriginURL()
			assert.Equal(t, want, *have, give)
		}
	})

	t.Run(".SetOffline()", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		err := repo.Config.SetOffline(true)
		assert.NoError(t, err)
		offline, err := repo.Config.IsOffline()
		assert.Nil(t, err)
		assert.True(t, offline)
		err = repo.Config.SetOffline(false)
		assert.NoError(t, err)
		offline, err = repo.Config.IsOffline()
		assert.Nil(t, err)
		assert.False(t, offline)
	})
}
