package cli_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/shoenig/test"
)

func TestIndent(t *testing.T) {
	t.Parallel()

	t.Run("no indent", func(t *testing.T) {
		t.Parallel()
		have := cli.Indent("")
		test.EqOp(t, "  ", have)
	})

	t.Run("single line of text", func(t *testing.T) {
		t.Parallel()
		have := cli.Indent("hello")
		test.EqOp(t, "  hello", have)
	})

	t.Run("multi-line text", func(t *testing.T) {
		t.Parallel()
		have := cli.Indent("hello\nworld")
		test.EqOp(t, "  hello\n  world", have)
	})

	t.Run("multiple newlines", func(t *testing.T) {
		t.Parallel()
		give := "hello\n\nworld"
		have := cli.Indent(give)
		want := "  hello\n\n  world"
		test.EqOp(t, want, have)
	})
}
