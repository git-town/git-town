package math_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/gohacks/math"
	"github.com/shoenig/test/must"
)

func TestMax(t *testing.T) {
	t.Parallel()
	must.Eq(t, 2, math.Max(1, 2))
	must.Eq(t, 2, math.Max(2, 1))
}
