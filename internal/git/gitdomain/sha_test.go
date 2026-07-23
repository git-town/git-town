package gitdomain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/git-town/git-town/v24/pkg/asserts"
	"github.com/shoenig/test/must"
)

func TestSHA(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		sha := gitdomain.SHA("123456")
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
			sha := gitdomain.NewSHAOrPanic(stringss.Trim(text))
			must.EqOp(t, text, sha.String())
		})
		t.Run("does not allow empty values", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.NewSHAOrPanic("")
		})
		t.Run("does not allow spaces", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.NewSHAOrPanic("abc def")
		})
		t.Run("does not allow uppercase characters", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.NewSHAOrPanic("ABCDEF")
		})
		t.Run("does not allow non-hex characters", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.NewSHAOrPanic("abcdefg")
		})
	})

	t.Run("Truncate", func(t *testing.T) {
		t.Parallel()
		give := gitdomain.SHA("12345678901234567890123456789012")
		want := gitdomain.SHA("1234567")
		have := give.Truncate(7)
		must.EqOp(t, want, have)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"123456"`
		var have gitdomain.SHA
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := gitdomain.SHA("123456")
		must.EqOp(t, want, have)
	})
}
