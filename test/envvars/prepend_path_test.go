package envvars_test

import (
	"runtime"
	"testing"

	"github.com/git-town/git-town/v9/test/envvars"
	"github.com/stretchr/testify/assert"
)

func TestPrependPath(t *testing.T) {
	t.Parallel()
	t.Run("already contains the given path", func(t *testing.T) {
		t.Parallel()
		var give []string
		var want []string
		if runtime.GOOS == "windows" {
			give = []string{"ONE=1", "PATH=alpha;beta", "THREE=3"}
			want = []string{"ONE=1", "PATH=gamma;alpha;beta", "THREE=3"}
		} else {
			give = []string{"ONE=1", "PATH=alpha:beta", "THREE=3"}
			want = []string{"ONE=1", "PATH=gamma:alpha:beta", "THREE=3"}
		}
		have := envvars.PrependPath(give, "gamma")
		assert.Equal(t, want, have)
	})

	t.Run("does not contain the given path", func(t *testing.T) {
		t.Parallel()
		give := []string{"ONE=1", "TWO=2"}
		have := envvars.PrependPath(give, "alpha")
		want := []string{"ONE=1", "TWO=2", "PATH=alpha"}
		assert.Equal(t, want, have)
	})
}
