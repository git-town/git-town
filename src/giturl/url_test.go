package giturl_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/stretchr/testify/assert"
)

func TestHost(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"git@github.com:git-town/git-town.git":                 "github.com",
		"username@bitbucket.org:git-town/git-town.git":         "bitbucket.org",
		"https://github.com/git-town/git-town.git":             "github.com",
		"https://user@github.com/git-town/git-town.git":        "github.com",
		"https://user:secret@github.com/git-town/git-town.git": "github.com",
	}
	for give, want := range tests {
		have := giturl.Host(give)
		assert.Equal(t, want, have, give)
	}
}

func TestRepo(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"git@github.com:git-town/git-town.git":                 "git-town/git-town",
		"username@bitbucket.org:git-town/git-town.git":         "git-town/git-town",
		"https://github.com/git-town/git-town.git":             "git-town/git-town",
		"https://user@github.com/git-town/git-town.git":        "git-town/git-town",
		"https://user:secret@github.com/git-town/git-town.git": "git-town/git-town",
	}
	for give, want := range tests {
		have := giturl.Repo(give)
		assert.Equal(t, want, have, give)
	}
}
