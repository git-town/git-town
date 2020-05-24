package git_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunner_SetOffline(t *testing.T) {
	repo := CreateTestGitTownRepo(t)
	_ = repo.SetOffline(true)
	// assert.Nil(t, err)
	offline := repo.IsOffline()
	assert.True(t, offline)
	_ = repo.SetOffline(false)
	// assert.Nil(t, err)
	offline = repo.IsOffline()
	assert.False(t, offline)
}
