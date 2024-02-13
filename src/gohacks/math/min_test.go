package math_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/gohacks/math"
	"github.com/shoenig/test/must"
)

func TestMin(t *testing.T) {
	t.Parallel()
	must.Eq(t, 1, math.Min(1, 2))
	must.Eq(t, 1, math.Min(2, 1))
}
