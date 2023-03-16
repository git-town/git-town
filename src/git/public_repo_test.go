package git_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/stretchr/testify/assert"
)

func TestNewProdRepo(t *testing.T) {
	t.Parallel()
	debug := false
	repo := git.NewProdRepo(&debug)
	assert.Equal(t, repo.Config, repo.Config)
}

func TestPublicRepo(t *testing.T) {
}
