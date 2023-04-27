package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/fs"
	"github.com/stretchr/testify/assert"
)

func TestCreateFile(t *testing.T) {
	t.Parallel()
	t.Run("simple example", func(t *testing.T) {
		t.Parallel()
		runtime := commands.Create(t)
		err := fs.CreateFile(runtime.Dir(), "filename", "content")
		assert.Nil(t, err, "cannot create file in repo")
		content, err := os.ReadFile(filepath.Join(runtime.Dir(), "filename"))
		assert.Nil(t, err, "cannot read file")
		assert.Equal(t, "content", string(content))
	})

	t.Run("create file in subfolder", func(t *testing.T) {
		t.Parallel()
		runtime := commands.Create(t)
		err := fs.CreateFile(runtime.Dir(), "folder/filename", "content")
		assert.Nil(t, err, "cannot create file in repo")
		content, err := os.ReadFile(filepath.Join(runtime.Dir(), "folder/filename"))
		assert.Nil(t, err, "cannot read file")
		assert.Equal(t, "content", string(content))
	})
}
