package git_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/stretchr/testify/assert"
)

func TestNewProdRepo(t *testing.T) {
	t.Parallel()
	repo := git.NewProdRepo()
	assert.Equal(t, repo.Config, repo.Silent.Config)
}
