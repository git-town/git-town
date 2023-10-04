package cli_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/shoenig/test/must"
)

func TestIndent(t *testing.T) {
	t.Parallel()

	t.Run("no indent", func(t *testing.T) {
		t.Parallel()
		have := cli.Indent("")
		must.EqOp(t, "  ", have)
	})

	t.Run("single line of text", func(t *testing.T) {
		t.Parallel()
		have := cli.Indent("hello")
		must.EqOp(t, "  hello", have)
	})

	t.Run("multi-line text", func(t *testing.T) {
		t.Parallel()
		have := cli.Indent("hello\nworld")
		must.EqOp(t, "  hello\n  world", have)
	})

	t.Run("multiple newlines", func(t *testing.T) {
		t.Parallel()
		give := "hello\n\nworld"
		have := cli.Indent(give)
		want := "  hello\n\n  world"
		must.EqOp(t, want, have)
	})
}
