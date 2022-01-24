package cli_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/stretchr/testify/assert"
)

func TestIndent_empty(t *testing.T) {
	t.Parallel()
	have := cli.Indent("")
	assert.Equal(t, have, "  ")
}

func TestIndent_singleLine(t *testing.T) {
	t.Parallel()
	have := cli.Indent("hello")
	assert.Equal(t, have, "  hello")
}

func TestIndent_multiLine(t *testing.T) {
	t.Parallel()
	have := cli.Indent("hello\nworld")
	assert.Equal(t, have, "  hello\n  world")
}

func TestIndent_multipleNewlines(t *testing.T) {
	t.Parallel()
	give := "hello\n\nworld"
	want := "  hello\n\n  world"
	have := cli.Indent(give)
	assert.Equal(t, want, have)
}
