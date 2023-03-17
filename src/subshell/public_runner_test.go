package subshell_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	tests := map[string][]string{
		"[branch] git checkout foo": {"branch", "git", "checkout", "foo"},
	}
	for want, give := range tests {
		have := subshell.FormatCommand(give[0], give[1], give[2:]...)
		assert.Equal(t, want, have)
	}
}
