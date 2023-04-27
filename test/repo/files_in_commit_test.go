package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/fs"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestFilesInCommit(t *testing.T) {
	t.Parallel()
	runtime := repo.Create(t)
	err := fs.CreateFile(runtime.Dir(), "f1.txt", "one")
	assert.NoError(t, err)
	err = fs.CreateFile(runtime.Dir(), "f2.txt", "two")
	assert.NoError(t, err)
	err = repo.StageFiles(&runtime, "f1.txt", "f2.txt")
	assert.NoError(t, err)
	err = repo.CommitStagedChanges(&runtime, "stuff")
	assert.NoError(t, err)
	commits, err := repo.Commits(&runtime, []string{}, "initial")
	assert.NoError(t, err)
	assert.Len(t, commits, 1)
	fileNames, err := repo.FilesInCommit(&runtime, commits[0].SHA)
	assert.NoError(t, err)
	assert.Equal(t, []string{"f1.txt", "f2.txt"}, fileNames)
}
