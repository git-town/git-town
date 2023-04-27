package commands_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/commands"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	t.Parallel()
	repo := commands.Create(t)
	origin := commands.Create(t)
	err := commands.AddRemote(&repo, config.OriginRemote, origin.Dir())
	assert.NoError(t, err)
	err = commands.Fetch(&repo)
	assert.NoError(t, err)
}
