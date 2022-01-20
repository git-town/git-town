package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/hosting/helpers"
	"github.com/stretchr/testify/assert"
)

func TestURLHostname(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"git@github.com:git-town/git-town.git":                 "github.com",
		"https://github.com/git-town/git-town.git":             "github.com",
		"https://user:secret@github.com/git-town/git-town.git": "github.com",
	}
	for give, want := range tests {
		have := helpers.URLHostname(give)
		assert.Equal(t, want, have)
	}
}
