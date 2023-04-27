package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/git"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommit(t *testing.T) {
	t.Parallel()
	t.Run("minimal arguments", func(t *testing.T) {
		t.Parallel()
		runtime := repo.Create(t)
		err := repo.CreateCommit(&runtime, git.Commit{
			Branch:      "initial",
			FileName:    "hello.txt",
			FileContent: "hello world",
			Message:     "test commit",
		})
		assert.NoError(t, err)
		commits, err := repo.Commits(&runtime, []string{"FILE NAME", "FILE CONTENT"}, "initial")
		assert.NoError(t, err)
		assert.Len(t, commits, 1)
		assert.Equal(t, "hello.txt", commits[0].FileName)
		assert.Equal(t, "hello world", commits[0].FileContent)
		assert.Equal(t, "test commit", commits[0].Message)
		assert.Equal(t, "initial", commits[0].Branch)
	})

	t.Run("set the author", func(t *testing.T) {
		t.Parallel()
		runtime := repo.Create(t)
		err := repo.CreateCommit(&runtime, git.Commit{
			Branch:      "initial",
			FileName:    "hello.txt",
			FileContent: "hello world",
			Message:     "test commit",
			Author:      "developer <developer@example.com>",
		})
		assert.NoError(t, err)
		commits, err := repo.Commits(&runtime, []string{"FILE NAME", "FILE CONTENT"}, "initial")
		assert.NoError(t, err)
		assert.Len(t, commits, 1)
		assert.Equal(t, "hello.txt", commits[0].FileName)
		assert.Equal(t, "hello world", commits[0].FileContent)
		assert.Equal(t, "test commit", commits[0].Message)
		assert.Equal(t, "initial", commits[0].Branch)
		assert.Equal(t, "developer <developer@example.com>", commits[0].Author)
	})
}
