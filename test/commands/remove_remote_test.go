package commands_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/commands"
	"github.com/stretchr/testify/assert"
)

func TestRemoveRemote(t *testing.T) {
	t.Parallel()
	repo := commands.Create(t)
	origin := commands.Create(t)
	err := commands.AddRemote(&repo, config.OriginRemote, origin.Dir())
	assert.NoError(t, err)
	err = commands.RemoveRemote(&repo, config.OriginRemote)
	assert.NoError(t, err)
	remotes, err := repo.Remotes()
	assert.NoError(t, err)
	assert.Len(t, remotes, 0)
}
