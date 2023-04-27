package commands_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/commands"
	"github.com/stretchr/testify/assert"
)

func TestPushBranchToRemote(t *testing.T) {
	t.Parallel()
	dev := commands.Create(t)
	origin := commands.Create(t)
	err := commands.AddRemote(&dev, config.OriginRemote, origin.Dir())
	assert.NoError(t, err)
	err = commands.CreateBranch(&dev, "b1", "initial")
	assert.NoError(t, err)
	err = commands.PushBranchToRemote(&dev, "b1", config.OriginRemote)
	assert.NoError(t, err)
	branches, err := origin.LocalBranchesMainFirst("initial")
	assert.NoError(t, err)
	assert.Equal(t, []string{"initial", "b1"}, branches)
}
