package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestStringify(t *testing.T) {
	t.Parallel()

	t.Run("0 elements", func(t *testing.T) {
		t.Parallel()
		have := slice.Stringify([]configdomain.BranchType{})
		want := []string{}
		must.Eq(t, want, have)
	})

	t.Run("1 element", func(t *testing.T) {
		t.Parallel()
		have := slice.Stringify([]configdomain.BranchType{configdomain.BranchTypeMainBranch})
		want := []string{"main"}
		must.Eq(t, want, have)
	})

	t.Run("2 elements", func(t *testing.T) {
		t.Parallel()
		have := slice.Stringify([]configdomain.BranchType{configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch})
		want := []string{"main", "perennial"}
		must.Eq(t, want, have)
	})
}
