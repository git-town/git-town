package git_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunner_SetOffline(t *testing.T) {
	repo := CreateTestGitTownRepo(t)
	err := repo.SetOffline(true)
	assert.NoError(t, err)
	offline := repo.IsOffline()
	assert.True(t, offline)
	err = repo.SetOffline(false)
	assert.NoError(t, err)
	offline = repo.IsOffline()
	assert.False(t, offline)
}
