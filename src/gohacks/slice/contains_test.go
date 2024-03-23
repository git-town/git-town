package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestContains(t *testing.T) {
	t.Parallel()
	give := []string{"one", "two"}
	must.True(t, slice.Contains(give, "one"))
	must.True(t, slice.Contains(give, "two"))
	must.False(t, slice.Contains(give, "three"))
}
