package git_test

import (
	"testing"

	"github.com/git-town/git-town/v14/test/git"
	"github.com/shoenig/test/must"
)

func TestLocations(t *testing.T) {
	t.Parallel()

	t.Run("NewLocations", func(t *testing.T) {
		t.Parallel()
		tests := map[string]git.Locations{
			"local":                   {git.LocationLocal},
			"origin":                  {git.LocationOrigin},
			"upstream":                {git.LocationUpstream},
			"local, origin":           {git.LocationLocal, git.LocationOrigin},
			"local, upstream":         {git.LocationLocal, git.LocationUpstream},
			"local, origin, upstream": {git.LocationLocal, git.LocationOrigin, git.LocationUpstream},
		}
		for give, want := range tests {
			have := git.NewLocations(give)
			must.Eq(t, want, have)
		}
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		t.Run("has the element", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationLocal, git.LocationOrigin}
			have := locations.Contains(git.LocationOrigin)
			must.True(t, have)
		})
		t.Run("does not have the element", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationLocal}
			have := locations.Contains(git.LocationOrigin)
			must.False(t, have)
		})
	})

	t.Run("Is", func(t *testing.T) {
		t.Parallel()
		t.Run("match with one element", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationOrigin}
			must.True(t, locations.Is(git.LocationOrigin))
		})
		t.Run("match with multiple elements", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationLocal, git.LocationOrigin}
			must.True(t, locations.Is(git.LocationLocal, git.LocationOrigin))
		})
		t.Run("wrong type", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationOrigin}
			must.False(t, locations.Is(git.LocationLocal))
		})
		t.Run("contains more elements", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationLocal, git.LocationOrigin}
			must.False(t, locations.Is(git.LocationLocal))
		})
		t.Run("contains fewer elements", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationLocal}
			must.False(t, locations.Is(git.LocationLocal, git.LocationOrigin))
		})
	})

	t.Run("Matches", func(t *testing.T) {
		t.Parallel()
		t.Run("has exactly the given elements", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationLocal, git.LocationOrigin}
			have := locations.Matches(git.LocationOrigin, git.LocationLocal)
			must.True(t, have)
		})
		t.Run("has the given elements and more", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationLocal, git.LocationOrigin, git.LocationCoworker}
			have := locations.Matches(git.LocationOrigin, git.LocationLocal)
			must.False(t, have)
		})
		t.Run("has not all of the given elements", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationLocal}
			have := locations.Matches(git.LocationOrigin, git.LocationLocal)
			must.False(t, have)
		})
		t.Run("has other elements", func(t *testing.T) {
			t.Parallel()
			locations := git.Locations{git.LocationLocal}
			have := locations.Matches(git.LocationOrigin)
			must.False(t, have)
		})
	})
}
