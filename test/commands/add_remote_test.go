package commands_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/runtime"
	"github.com/stretchr/testify/assert"
)

func TestAddRemote(t *testing.T) {
	t.Parallel()
	t.Run("remote doesn't exist", func(t *testing.T) {
		t.Parallel()
		dev := runtime.Create(t)
		remotes, err := dev.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{}, remotes)
		origin := runtime.Create(t)
		err = commands.AddRemote(&dev, config.OriginRemote, origin.Dir())
		assert.NoError(t, err)
		remotes, err = dev.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin"}, remotes)
	})
}
