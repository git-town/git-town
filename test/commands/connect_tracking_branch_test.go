package commands_test

import (
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/fs"
	"github.com/git-town/git-town/v8/test/runtime"
	"github.com/stretchr/testify/assert"
)

func TestConnectTrackingBranch(t *testing.T) {
	t.Parallel()
	// replicating the situation this is used in,
	// connecting branches of repos with the same commits in them
	origin := runtime.Create(t)
	repoDir := filepath.Join(t.TempDir(), "repo") // need a non-existing directory
	err := fs.CopyDirectory(origin.Dir(), repoDir)
	assert.NoError(t, err)
	runtime := runtime.New(repoDir, repoDir, "")
	err = commands.AddRemote(&runtime, config.OriginRemote, origin.Dir())
	assert.NoError(t, err)
	err = commands.Fetch(&runtime)
	assert.NoError(t, err)
	err = commands.ConnectTrackingBranch(&runtime, "initial")
	assert.NoError(t, err)
	err = commands.PushBranch(&runtime)
	assert.NoError(t, err)
}
