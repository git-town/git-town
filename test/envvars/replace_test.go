package envvars_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/envvars"
	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	t.Run("contains the given key", func(t *testing.T) {
		t.Parallel()
		give := []string{"ONE=1", "TWO=2", "THREE=3"}
		have := envvars.Replace(give, "TWO", "another")
		want := []string{"ONE=1", "TWO=another", "THREE=3"}
		assert.Equal(t, have, want)
	})

	t.Run("doesn't contain the given key", func(t *testing.T) {
		t.Parallel()
		give := []string{"ONE=1", "TWO=2"}
		have := envvars.Replace(give, "THREE", "new")
		want := []string{"ONE=1", "TWO=2", "THREE=new"}
		assert.Equal(t, have, want)
	})
}
