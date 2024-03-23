package format_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/cli/format"
	"github.com/shoenig/test/must"
)

func TestIndent(t *testing.T) {
	t.Parallel()

	t.Run("no indent", func(t *testing.T) {
		t.Parallel()
		have := format.Indent("")
		must.EqOp(t, "  ", have)
	})

	t.Run("single line of text", func(t *testing.T) {
		t.Parallel()
		have := format.Indent("hello")
		must.EqOp(t, "  hello", have)
	})

	t.Run("multi-line text", func(t *testing.T) {
		t.Parallel()
		have := format.Indent("hello\nworld")
		must.EqOp(t, "  hello\n  world", have)
	})

	t.Run("multiple newlines", func(t *testing.T) {
		t.Parallel()
		give := "hello\n\nworld"
		have := format.Indent(give)
		want := "  hello\n\n  world"
		must.EqOp(t, want, have)
	})
}
