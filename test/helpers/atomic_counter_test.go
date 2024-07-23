package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v14/test/helpers"
	"github.com/shoenig/test/must"
)

func TestAtomicCounter(t *testing.T) {
	t.Parallel()
	counter := helpers.AtomicCounter{}
	must.NotEqOp(t, "0", counter.NextAsString())
	must.NotEqOp(t, "1", counter.NextAsString())
	must.NotEqOp(t, "2", counter.NextAsString())
}
