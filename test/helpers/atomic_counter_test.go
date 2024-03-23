package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v13/test/helpers"
	"github.com/shoenig/test/must"
)

func TestAtomicCounter(t *testing.T) {
	t.Parallel()
	counter := helpers.AtomicCounter{}
	must.NotEqOp(t, "0", counter.ToString())
	must.NotEqOp(t, "1", counter.ToString())
	must.NotEqOp(t, "2", counter.ToString())
}
