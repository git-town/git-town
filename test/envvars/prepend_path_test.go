package envvars_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/envvars"
	"github.com/stretchr/testify/assert"
)

func TestPrependPath(t *testing.T) {
	t.Parallel()
	t.Run("already contains the given path", func(t *testing.T) {
		t.Parallel()
		give := []string{"ONE=1", "PATH=alpha:beta", "THREE=3"}
		have := envvars.PrependPath(give, "gamma")
		want := []string{"ONE=1", "PATH=gamma:alpha:beta", "THREE=3"}
		assert.Equal(t, have, want)
	})

	t.Run("does not contain the given path", func(t *testing.T) {
		t.Parallel()
		give := []string{"ONE=1", "TWO=2"}
		have := envvars.PrependPath(give, "alpha")
		want := []string{"ONE=1", "TWO=2", "PATH=alpha"}
		assert.Equal(t, have, want)
	})
}
