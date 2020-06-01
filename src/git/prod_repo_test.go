package git_test

import (
	"testing"

	"github.com/git-town/git-town/src/git"
	"github.com/stretchr/testify/assert"
)

func TestNewProdRepo(t *testing.T) {
	repo := git.NewProdRepo()
	assert.Equal(t, repo.ConfigurationInterface, repo.Silent.ConfigurationInterface)
}
