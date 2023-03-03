package config_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/test"
	"github.com/stretchr/testify/assert"
)

func TestGitTown(t *testing.T) {
	t.Parallel()
	t.Run(".SetOffline()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateTestGitTownRepo(t)
		err := repo.Config.SetOffline(true)
		assert.NoError(t, err)
		offline, err := repo.Config.IsOffline()
		assert.Nil(t, err)
		assert.True(t, offline)
		err = repo.Config.SetOffline(false)
		assert.NoError(t, err)
		offline, err = repo.Config.IsOffline()
		assert.Nil(t, err)
		assert.False(t, offline)
	})
}

func TestToAliasType(t *testing.T) {
	t.Parallel()
	tests := map[string]config.AliasType{
		"append":           config.AliasTypeAppend,
		"diff-parent":      config.AliasTypeDiffParent,
		"hack":             config.AliasTypeHack,
		"kill":             config.AliasTypeKill,
		"new-pull-request": config.AliasTypeNewPullRequest,
		"prepend":          config.AliasTypePrepend,
		"prune-branches":   config.AliasTypePruneBranches,
		"rename-branch":    config.AliasTypeRenameBranch,
		"repo":             config.AliasTypeRepo,
		"ship":             config.AliasTypeShip,
		"sync":             config.AliasTypeSync,
	}
	for give, want := range tests {
		have, err := config.ToAliasType(give)
		assert.Nil(t, err)
		assert.Equal(t, want, have)
	}
}
