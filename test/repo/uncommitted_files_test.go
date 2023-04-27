package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/fs"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestUncommittedFiles(t *testing.T) {
	t.Parallel()
	runtime := repo.Create(t)
	err := fs.CreateFile(runtime.Dir(), "f1.txt", "one")
	assert.NoError(t, err)
	err = fs.CreateFile(runtime.Dir(), "f2.txt", "two")
	assert.NoError(t, err)
	files, err := repo.UncommittedFiles(&runtime)
	assert.NoError(t, err)
	assert.Equal(t, []string{"f1.txt", "f2.txt"}, files)
}
