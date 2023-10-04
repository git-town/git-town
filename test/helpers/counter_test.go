package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v9/test/helpers"
	"github.com/shoenig/test/must"
)

func TestCounter(t *testing.T) {
	t.Parallel()
	counter := helpers.Counter{}
	must.NotEqOp(t, "0", counter.ToString())
	must.NotEqOp(t, "1", counter.ToString())
	must.NotEqOp(t, "2", counter.ToString())
}
