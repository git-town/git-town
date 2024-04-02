package git_test

import (
	"testing"

	testgit "github.com/git-town/git-town/v13/test/git"
	"github.com/shoenig/test/must"
)

func TestNewLocations(t *testing.T) {
	t.Parallel()
	tests := map[string]testgit.Locations{
		"local":                   {testgit.LocationLocal},
		"origin":                  {testgit.LocationOrigin},
		"upstream":                {testgit.LocationUpstream},
		"local, origin":           {testgit.LocationLocal, testgit.LocationOrigin},
		"local, upstream":         {testgit.LocationLocal, testgit.LocationUpstream},
		"local, origin, upstream": {testgit.LocationLocal, testgit.LocationOrigin, testgit.LocationUpstream},
	}
	for give, want := range tests {
		have := testgit.NewLocations(give)
		must.Eq(t, want, have)
	}
}
