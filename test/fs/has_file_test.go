package fs_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/fs"
	"github.com/stretchr/testify/assert"
)

func TestHasFile(t *testing.T) {
	t.Parallel()
	runtime := commands.Create(t)
	err := fs.CreateFile(runtime.Dir(), "f1.txt", "one")
	assert.NoError(t, err)
	has, err := fs.HasFile(runtime.WorkingDir, "f1.txt", "one")
	assert.NoError(t, err)
	assert.True(t, has)
	_, err = fs.HasFile(runtime.WorkingDir, "f1.txt", "zonk")
	assert.Error(t, err)
	_, err = fs.HasFile(runtime.WorkingDir, "zonk.txt", "one")
	assert.Error(t, err)
}
