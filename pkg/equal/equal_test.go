package equal_test

import (
	"testing"

	"github.com/git-town/git-town/v22/pkg/equal"
	"github.com/shoenig/test/must"
)

func TestEqual(t *testing.T) {
	t.Parallel()
	must.True(t, equal.Equal(1, 1))
	must.False(t, equal.Equal(1, 2))
}
