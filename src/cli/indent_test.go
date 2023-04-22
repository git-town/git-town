package cli_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/cli"
	"github.com/stretchr/testify/assert"
)

func TestIndent(t *testing.T) {
	t.Parallel()
	t.Run("no indent", func(t *testing.T) {
		t.Parallel()
		have := cli.Indent("")
		assert.Equal(t, have, "  ")
	})

	t.Run("single line of text", func(t *testing.T) {
		t.Parallel()
		have := cli.Indent("hello")
		assert.Equal(t, have, "  hello")
	})

	t.Run("multi-line text", func(t *testing.T) {
		t.Parallel()
		have := cli.Indent("hello\nworld")
		assert.Equal(t, have, "  hello\n  world")
	})

	t.Run("multiple newlines", func(t *testing.T) {
		t.Parallel()
		give := "hello\n\nworld"
		want := "  hello\n\n  world"
		have := cli.Indent(give)
		assert.Equal(t, want, have)
	})
}
