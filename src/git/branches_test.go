package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/stretchr/testify/assert"
)

func TestAncestry(t *testing.T) {
	t.Parallel()

	t.Run("Ancestors", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all ancestor branches, oldest first", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{
				git.Branch{
					Name:   "three",
					Parent: "two",
				},
				git.Branch{
					Name:   "two",
					Parent: "one",
				},
				git.Branch{
					Name:   "one",
					Parent: "main",
				},
				git.Branch{
					Name:   "main",
					Parent: "",
				},
				git.Branch{
					Name:   "other",
					Parent: "one",
				},
			}
			want := []string{"main", "one", "two", "three"}
			have := bs.Ancestors("three").BranchNames()
			assert.Equal(t, want, have)
		})
		t.Run("no ancestors", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{}
			have := bs.Ancestors("foo")
			want := git.Branches{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the names of all branches in this collection, ordered the same as the collection", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{
				git.Branch{Name: "one"},
				git.Branch{Name: "two"},
				git.Branch{Name: "three"},
			}
			want := []string{"one", "two", "three"}
			have := bs.BranchNames()
			assert.Equal(t, want, have)
		})
	})

	t.Run("Children", func(t *testing.T) {
		t.Parallel()
		t.Run("provides all children of the given branch, ordered alphabetically", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{
				git.Branch{
					Name:   "alpha",
					Parent: "",
				},
				git.Branch{
					Name:   "beta2",
					Parent: "alpha",
				},
				git.Branch{
					Name:   "beta1",
					Parent: "alpha",
				},
			}
			have := bs.Children("alpha").BranchNames()
			want := []string{"beta1", "beta2"}
			assert.Equal(t, want, have)
		})
		t.Run("provides only the immediate children, i.e. no grandchildren", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{
				git.Branch{
					Name:   "one",
					Parent: "",
				},
				git.Branch{
					Name:   "two",
					Parent: "one",
				},
				git.Branch{
					Name:   "three",
					Parent: "two",
				},
			}
			have := bs.Children("one").BranchNames()
			want := []string{"two"}
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{}
			have := bs.Children("alpha")
			want := []string{}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		bs := git.Branches{
			git.Branch{Name: "one"},
			git.Branch{Name: "two"},
		}
		assert.True(t, bs.Contains("one"))
		assert.True(t, bs.Contains("two"))
		assert.False(t, bs.Contains("zonk"))
	})

	t.Run("HasParent", func(t *testing.T) {
		t.Parallel()
		bs := git.Branches{
			git.Branch{
				Name:   "one",
				Parent: "",
			},
			git.Branch{
				Name:   "two",
				Parent: "one",
			},
		}
		assert.True(t, bs.HasParent("two"))
		assert.False(t, bs.HasParent("one"))
	})

	t.Run("IsAncestor", func(t *testing.T) {
		t.Run("recognizes greatgrandparent", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{
				git.Branch{
					Name:   "one",
					Parent: "",
				},
				git.Branch{
					Name:   "two",
					Parent: "one",
				},
				git.Branch{
					Name:   "three",
					Parent: "two",
				},
				git.Branch{
					Name:   "four",
					Parent: "three",
				},
			}
			assert.True(t, bs.IsAncestor("one", "four"))
		})
		t.Run("child branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{
				git.Branch{
					Name:   "one",
					Parent: "",
				},
				git.Branch{
					Name:   "two",
					Parent: "one",
				},
			}
			assert.False(t, bs.IsAncestor("two", "one"))
		})
		t.Run("unrelated branches are not ancestors", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{
				git.Branch{
					Name:   "one",
					Parent: "",
				},
				git.Branch{
					Name:   "two",
					Parent: "one",
				},
				git.Branch{
					Name:   "two",
					Parent: "one",
				},
			}
			assert.False(t, bs.IsAncestor("two", "three"))
		})
	})

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := git.Branches{
			git.Branch{
				Name:       "up-to-date",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.Branch{
				Name:       "ahead",
				SyncStatus: git.SyncStatusAhead,
			},
			git.Branch{
				Name:       "behind",
				SyncStatus: git.SyncStatusBehind,
			},
			git.Branch{
				Name:       "local-only",
				SyncStatus: git.SyncStatusLocalOnly,
			},
			git.Branch{
				Name:       "remote-only",
				SyncStatus: git.SyncStatusRemoteOnly,
			},
			git.Branch{
				Name:       "deleted-at-remote",
				SyncStatus: git.SyncStatusDeletedAtRemote,
			},
		}
		want := []string{"up-to-date", "ahead", "behind", "local-only"}
		have := bs.LocalBranches().BranchNames()
		assert.Equal(t, want, have)
	})

	t.Run("Lookup", func(t *testing.T) {
		t.Parallel()
		bs := git.Branches{
			git.Branch{
				Name: "one",
			},
			git.Branch{
				Name: "two",
			},
		}
		assert.Equal(t, "one", bs.Lookup("one").Name)
		assert.Equal(t, "two", bs.Lookup("two").Name)
		assert.Nil(t, bs.Lookup("zonk"))
	})

	t.Run("OrderedHierarchically", func(t *testing.T) {
		t.Parallel()
		bs := git.Branches{
			git.Branch{
				Name:   "main",
				Parent: "",
			},
			git.Branch{
				Name:   "1",
				Parent: "main",
			},
			git.Branch{
				Name:   "1A",
				Parent: "1",
			},
			git.Branch{
				Name:   "1B",
				Parent: "one",
			},
			git.Branch{
				Name:   "1A1",
				Parent: "1A",
			},
			git.Branch{
				Name:   "1A2",
				Parent: "1A",
			},
			git.Branch{
				Name:   "2",
				Parent: "main",
			},
		}
		want := []string{"main", "1", "2", "1A", "1B", "1A1", "1A2"}
		have := bs.OrderedHierarchically().BranchNames()
		assert.Equal(t, want, have)
	})

	t.Run("Roots", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple roots with nested child branches", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{
				git.Branch{
					Name:   "main",
					Parent: "",
				},
				git.Branch{
					Name:   "one",
					Parent: "main",
				},
				git.Branch{
					Name:   "two",
					Parent: "one",
				},
				git.Branch{
					Name:   "alpha",
					Parent: "main",
				},
				git.Branch{
					Name:   "beta",
					Parent: "alpha",
				},
				git.Branch{
					Name:   "prod",
					Parent: "",
				},
				git.Branch{
					Name:   "hotfix1",
					Parent: "prod",
				},
			}
			want := []string{"main", "prod"}
			have := bs.Roots().BranchNames()
			assert.Equal(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			bs := git.Branches{}
			have := bs.Roots()
			want := []string{}
			assert.Equal(t, want, have)
		})
	})
}
