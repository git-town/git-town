package gitdomain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/pkg/asserts"
	"github.com/shoenig/test/must"
)

func TestSHA(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		sha := gitdomain.NewSHA("123456")
		have, err := json.MarshalIndent(sha, "", "  ")
		must.NoError(t, err)
		want := `"123456"`
		must.EqOp(t, want, string(have))
	})

	t.Run("NewSHA and String", func(t *testing.T) {
		t.Parallel()
		t.Run("allows lowercase hex characters", func(t *testing.T) {
			t.Parallel()
			text := "1234567890abcdef"
			sha := gitdomain.NewSHA(text)
			must.EqOp(t, text, sha.String())
		})
		t.Run("does not allow empty values", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.NewSHA("")
		})
		t.Run("does not allow spaces", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.NewSHA("abc def")
		})
		t.Run("does not allow uppercase characters", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.NewSHA("ABCDEF")
		})
		t.Run("does not allow non-hex characters", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.NewSHA("abcdefg")
		})
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"123456"`
		var have gitdomain.SHA
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := gitdomain.NewSHA("123456")
		must.EqOp(t, want, have)
	})
}
